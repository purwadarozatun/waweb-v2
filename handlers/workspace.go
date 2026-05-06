package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"waweb-v2/config"
	"waweb-v2/database"
	"waweb-v2/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func NewWorkspace(c *fiber.Ctx) error {
	return c.Render("new-workspace", fiber.Map{
		"Title": "New Workspace",
	})
}

func CreateWorkspace(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	name := c.FormValue("name")

	workspace := &models.Workspace{
		ID:          uuid.New().String(),
		UserID:      userID,
		Name:        name,
		IsPublished: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := database.CreateWorkspace(workspace); err != nil {
		return c.Status(500).SendString("Error creating workspace")
	}

	// Redirect to template picker
	c.Set("HX-Redirect", fmt.Sprintf("/templates?workspace=%s", workspace.ID))
	return c.SendStatus(200)
}

func WorkspaceDetail(c *fiber.Ctx) error {
	workspaceID := c.Params("id")
	userID := c.Locals("userID").(string)

	workspace, exists := database.GetWorkspaceByID(workspaceID)
	if !exists || workspace.UserID != userID {
		return c.Status(404).SendString("Workspace not found")
	}

	return c.JSON(workspace)
}

func DeleteWorkspace(c *fiber.Ctx) error {
	workspaceID := c.Params("id")
	userID := c.Locals("userID").(string)

	workspace, exists := database.GetWorkspaceByID(workspaceID)
	if !exists || workspace.UserID != userID {
		return c.Status(404).SendString("Workspace not found")
	}

	// Delete published files if exists
	if workspace.IsPublished && workspace.Subdomain != "" {
		cfg := config.LoadConfig()
		sitePath := filepath.Join(cfg.HostFolder, workspace.Subdomain)
		os.RemoveAll(sitePath)
	}

	if err := database.DeleteWorkspace(workspaceID); err != nil {
		return c.Status(500).SendString("Error deleting workspace")
	}

	return c.SendStatus(200)
}

func ListTemplates(c *fiber.Ctx) error {
	workspaceID := c.Query("workspace")

	// Redirect to new workspace if no workspace ID provided
	if workspaceID == "" {
		return c.Redirect("/workspace/new")
	}

	cfg := config.LoadConfig()
	templatesPath := cfg.TemplatesFolder

	var templates []*models.Template

	// Read templates directory
	entries, err := os.ReadDir(templatesPath)
	if err != nil {
		return c.Render("templates", fiber.Map{
			"Title":       "Select Template",
			"WorkspaceID": workspaceID,
			"Templates":   templates,
		})
	}

	for _, entry := range entries {
		if entry.IsDir() {
			specPath := filepath.Join(templatesPath, entry.Name(), "spec.json")
			if data, err := os.ReadFile(specPath); err == nil {
				var spec map[string]interface{}
				if err := json.Unmarshal(data, &spec); err == nil {
					template := &models.Template{
						ID:          entry.Name(),
						Name:        getStringFromMap(spec, "name", entry.Name()),
						Description: getStringFromMap(spec, "description", ""),
						Thumbnail:   getStringFromMap(spec, "thumbnail", "/static/default-template.png"),
						Spec:        spec,
					}
					templates = append(templates, template)
				}
			}
		}
	}

	return c.Render("templates", fiber.Map{
		"Title":       "Select Template",
		"WorkspaceID": workspaceID,
		"Templates":   templates,
	})
}

func GetTemplate(c *fiber.Ctx) error {
	templateID := c.Params("id")
	workspaceID := c.Query("workspace")
	userID := c.Locals("userID").(string)

	fmt.Printf("[GetTemplate] templateID=%s, workspaceID=%s, userID=%s\n", templateID, workspaceID, userID)

	if workspaceID == "" {
		return c.Status(400).SendString("Missing workspace ID")
	}

	// Get workspace
	workspace, exists := database.GetWorkspaceByID(workspaceID)
	if !exists || workspace.UserID != userID {
		return c.Status(404).SendString("Workspace not found")
	}

	// Update workspace with template
	workspace.TemplateID = templateID
	workspace.UpdatedAt = time.Now()

	fmt.Printf("[GetTemplate] Setting template ID to: %s for workspace: %s\n", templateID, workspaceID)

	// Load template spec as initial values
	cfg := config.LoadConfig()
	specPath := filepath.Join(cfg.TemplatesFolder, templateID, "spec.json")
	if data, err := os.ReadFile(specPath); err == nil {
		var spec map[string]interface{}
		if err := json.Unmarshal(data, &spec); err == nil {
			// Extract default values from spec
			if fields, ok := spec["fields"].(map[string]interface{}); ok {
				values := make(map[string]interface{})
				for key, field := range fields {
					if fieldMap, ok := field.(map[string]interface{}); ok {
						if defaultVal, exists := fieldMap["default"]; exists {
							values[key] = defaultVal
						}
					}
				}
				workspace.Values = values
			}
		}
	}

	database.UpdateWorkspace(workspace)

	// Redirect to builder
	c.Set("HX-Redirect", fmt.Sprintf("/builder/%s", workspaceID))
	return c.SendStatus(200)
}

func getStringFromMap(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return defaultValue
}
