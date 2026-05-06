package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"waweb-v2/config"
	"waweb-v2/database"

	"github.com/gofiber/fiber/v2"
)

// PreviewSite serves live preview of workspace before publishing
func PreviewSite(c *fiber.Ctx) error {
	workspaceID := c.Params("id")
	requestedFile := c.Params("*")

	// Clean up requested file path
	requestedFile = strings.TrimPrefix(requestedFile, "/")

	// Default to index.html
	if requestedFile == "" {
		requestedFile = "index.html"
	}

	// Get workspace
	workspace, exists := database.GetWorkspaceByID(workspaceID)
	if !exists {
		return c.Status(404).SendString("Workspace not found")
	}

	if workspace.TemplateID == "" {
		return c.Status(404).SendString("No template selected")
	}

	cfg := config.LoadConfig()
	templatePath := filepath.Join(cfg.TemplatesFolder, workspace.TemplateID)

	// Special handling for values.json - generate from workspace.Values
	if requestedFile == "values.json" || requestedFile == "assets/values.json" || requestedFile == "assets/value.json" {
		valuesJSON, err := json.Marshal(workspace.Values)
		if err != nil {
			return c.Status(500).SendString("Error generating values.json")
		}

		c.Set("Content-Type", "application/json")
		c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		return c.Send(valuesJSON)
	}

	// Serve static files from template folder
	filePath := filepath.Join(templatePath, requestedFile)

	// Security check - prevent directory traversal
	absTemplatePath, _ := filepath.Abs(templatePath)
	absFilePath, _ := filepath.Abs(filePath)
	if !strings.HasPrefix(absFilePath, absTemplatePath) {
		return c.Status(403).SendString("Access denied")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(404).SendString("File not found: " + requestedFile)
	}

	// Set correct content type
	ext := filepath.Ext(requestedFile)
	switch ext {
	case ".html":
		c.Set("Content-Type", "text/html")
	case ".js":
		c.Set("Content-Type", "application/javascript")
	case ".css":
		c.Set("Content-Type", "text/css")
	case ".json":
		c.Set("Content-Type", "application/json")
	case ".png":
		c.Set("Content-Type", "image/png")
	case ".jpg", ".jpeg":
		c.Set("Content-Type", "image/jpeg")
	case ".svg":
		c.Set("Content-Type", "image/svg+xml")
	}

	return c.SendFile(filePath)
}

// RefreshPreview is a helper endpoint to trigger iframe reload via HTMX
func RefreshPreview(c *fiber.Ctx) error {
	workspaceID := c.Params("id")

	return c.SendString(fmt.Sprintf(`
		<script>
			const iframe = window.parent.document.getElementById('preview-iframe');
			if (iframe) {
				iframe.src = '/preview/%s/';
			}
		</script>
	`, workspaceID))
}
