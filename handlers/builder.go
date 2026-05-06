package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"waweb-v2/config"
	"waweb-v2/database"

	"github.com/gofiber/fiber/v2"
)

func Builder(c *fiber.Ctx) error {
	workspaceID := c.Params("id")
	userID := c.Locals("userID").(string)

	workspace, exists := database.GetWorkspaceByID(workspaceID)
	if !exists || workspace.UserID != userID {
		return c.Status(404).SendString("Workspace not found")
	}

	if workspace.TemplateID == "" {
		return c.Redirect("/templates?workspace=" + workspaceID)
	}

	fmt.Printf("[Builder] workspace=%s, templateID=%s\n", workspaceID, workspace.TemplateID)

	// Load template spec
	cfg := config.LoadConfig()
	specPath := filepath.Join(cfg.TemplatesFolder, workspace.TemplateID, "spec.json")
	fmt.Printf("[Builder] Loading spec from: %s\n", specPath)
	specData, err := os.ReadFile(specPath)
	if err != nil {
		fmt.Printf("Error reading spec file %s: %v\n", specPath, err)
		return c.Status(500).SendString("Error loading template: " + err.Error())
	}

	var spec map[string]interface{}
	if err := json.Unmarshal(specData, &spec); err != nil {
		fmt.Printf("Error parsing spec JSON: %v\n", err)
		return c.Status(500).SendString("Error parsing template spec: " + err.Error())
	}

	err = c.Render("builder", fiber.Map{
		"Title":     "Builder - " + workspace.Name,
		"Workspace": workspace,
		"Spec":      spec,
	})
	if err != nil {
		fmt.Printf("Error rendering builder template: %v\n", err)
		return c.Status(500).SendString("Error rendering template: " + err.Error())
	}
	return nil
}

func SaveBuilder(c *fiber.Ctx) error {
	workspaceID := c.Params("id")
	userID := c.Locals("userID").(string)

	workspace, exists := database.GetWorkspaceByID(workspaceID)
	if !exists || workspace.UserID != userID {
		return c.Status(404).SendString("Workspace not found")
	}

	// Merge incoming form fields into existing values.
	// This preserves image URLs saved by UploadAsset — we never wipe keys
	// that aren't present in the current form submission.
	if workspace.Values == nil {
		workspace.Values = make(map[string]interface{})
	}
	c.Request().PostArgs().VisitAll(func(key, value []byte) {
		workspace.Values[string(key)] = string(value)
	})
	workspace.UpdatedAt = time.Now()

	if err := database.UpdateWorkspace(workspace); err != nil {
		return c.Status(500).SendString("Error saving workspace")
	}

	return c.SendString(`
		<div class="alert alert-success" role="alert">
			Changes saved successfully!
		</div>
	`)
}

func PublishSite(c *fiber.Ctx) error {
	workspaceID := c.Params("id")
	userID := c.Locals("userID").(string)

	workspace, exists := database.GetWorkspaceByID(workspaceID)
	if !exists || workspace.UserID != userID {
		return c.Status(404).SendString("Workspace not found")
	}

	subdomain := c.FormValue("subdomain")
	if subdomain == "" {
		subdomain = strings.ToLower(strings.ReplaceAll(workspace.Name, " ", "-"))
	}

	// Validate subdomain
	if !isValidSubdomain(subdomain) {
		return c.Status(400).SendString("Invalid subdomain")
	}

	cfg := config.LoadConfig()

	// Create host folder if not exists
	if err := os.MkdirAll(cfg.HostFolder, 0755); err != nil {
		return c.Status(500).SendString("Error creating host folder")
	}

	// Copy template to host folder
	templatePath := filepath.Join(cfg.TemplatesFolder, workspace.TemplateID)
	sitePath := filepath.Join(cfg.HostFolder, subdomain)

	// Remove existing site if exists
	os.RemoveAll(sitePath)

	// Copy template
	if err := copyDir(templatePath, sitePath); err != nil {
		return c.Status(500).SendString("Error copying template: " + err.Error())
	}

	// Write values.json
	valuesData, err := json.MarshalIndent(workspace.Values, "", "  ")
	if err != nil {
		return c.Status(500).SendString("Error generating values")
	}

	if err := os.WriteFile(filepath.Join(sitePath, "values.json"), valuesData, 0644); err != nil {
		return c.Status(500).SendString("Error writing values")
	}

	// Update workspace
	workspace.Subdomain = subdomain
	workspace.IsPublished = true
	workspace.UpdatedAt = time.Now()

	if err := database.UpdateWorkspace(workspace); err != nil {
		return c.Status(500).SendString("Error updating workspace")
	}

	siteURL := fmt.Sprintf("/sites/%s/", subdomain)

	return c.SendString(fmt.Sprintf(`
		<div class="alert alert-success" role="alert">
			Site published successfully! 
			<a href="%s" target="_blank" class="alert-link">View Site</a>
		</div>
	`, siteURL))
}

func ServeSite(c *fiber.Ctx) error {
	subdomain := c.Params("subdomain")
	path := c.Params("*")

	cfg := config.LoadConfig()
	sitePath := filepath.Join(cfg.HostFolder, subdomain)

	// Check if site exists
	if _, err := os.Stat(sitePath); os.IsNotExist(err) {
		return c.Status(404).SendString("Site not found")
	}

	// Serve file
	if path == "" || path == "/" {
		path = "index.html"
	}

	filePath := filepath.Join(sitePath, path)

	// Security: prevent directory traversal
	if !strings.HasPrefix(filepath.Clean(filePath), sitePath) {
		return c.Status(403).SendString("Forbidden")
	}

	return c.SendFile(filePath)
}

// UploadAsset handles file uploads for image-type fields
func UploadAsset(c *fiber.Ctx) error {
	workspaceID := c.Params("id")
	userID := c.Locals("userID").(string)

	workspace, exists := database.GetWorkspaceByID(workspaceID)
	if !exists || workspace.UserID != userID {
		return c.Status(404).SendString("Workspace not found")
	}

	fieldKey := c.FormValue("field_key")
	if fieldKey == "" {
		return c.Status(400).SendString("Missing field_key")
	}

	// Validate field_key against template spec
	cfg := config.LoadConfig()
	specPath := filepath.Join(cfg.TemplatesFolder, workspace.TemplateID, "spec.json")
	specData, err := os.ReadFile(specPath)
	if err != nil {
		return c.Status(500).SendString("Error reading template spec")
	}
	var spec map[string]interface{}
	if err := json.Unmarshal(specData, &spec); err != nil {
		return c.Status(500).SendString("Error parsing template spec")
	}
	fields, _ := spec["fields"].(map[string]interface{})
	field, ok := fields[fieldKey]
	if !ok {
		return c.Status(400).SendString("Invalid field key")
	}
	fieldMap, _ := field.(map[string]interface{})
	fieldType, _ := fieldMap["type"].(string)
	if fieldType != "image" && fieldType != "file" {
		return c.Status(400).SendString("Field is not an image type")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).SendString("No file uploaded")
	}

	// Validate content type
	ct := file.Header.Get("Content-Type")
	allowedExt := map[string]string{
		"image/jpeg":    ".jpg",
		"image/png":     ".png",
		"image/gif":     ".gif",
		"image/webp":    ".webp",
		"image/svg+xml": ".svg",
	}
	ext, ok := allowedExt[ct]
	if !ok {
		return c.Status(400).SendString("Invalid file type. Allowed: JPEG, PNG, GIF, WebP, SVG")
	}

	// Validate file size (max 5 MB)
	if file.Size > 5*1024*1024 {
		return c.Status(400).SendString("File too large. Maximum 5 MB")
	}

	// Sanitize fieldKey for use in filesystem path
	safeKey := sanitizeKey(fieldKey)
	if safeKey == "" {
		return c.Status(400).SendString("Invalid field key characters")
	}

	uploadDir := filepath.Join("./uploads", workspaceID, safeKey)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.Status(500).SendString("Error creating upload directory")
	}

	safeFilename := safeKey + ext
	filePath := filepath.Join(uploadDir, safeFilename)

	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(500).SendString("Error saving file")
	}

	// Store absolute URL in workspace values
	fileURL := fmt.Sprintf("/uploads/%s/%s/%s", workspaceID, safeKey, safeFilename)
	if workspace.Values == nil {
		workspace.Values = make(map[string]interface{})
	}
	workspace.Values[fieldKey] = fileURL
	workspace.UpdatedAt = time.Now()

	if err := database.UpdateWorkspace(workspace); err != nil {
		return c.Status(500).SendString("Error updating workspace")
	}

	return c.JSON(fiber.Map{
		"url":       fileURL,
		"field_key": fieldKey,
	})
}

func sanitizeKey(key string) string {
	result := make([]byte, 0, len(key))
	for _, ch := range key {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '_' || ch == '-' {
			result = append(result, byte(ch))
		}
	}
	return string(result)
}

func isValidSubdomain(subdomain string) bool {
	if len(subdomain) == 0 || len(subdomain) > 63 {
		return false
	}
	for _, ch := range subdomain {
		if !((ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-') {
			return false
		}
	}
	return true
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip spec.json as it's not needed in published site
		if filepath.Base(path) == "spec.json" {
			return nil
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath)
	})
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
