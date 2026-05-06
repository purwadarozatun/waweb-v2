package handlers

import (
	"waweb-v2/database"

	"github.com/gofiber/fiber/v2"
)

func Dashboard(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	email := c.Locals("email").(string)

	workspaces := database.GetWorkspacesByUserID(userID)

	user, _ := database.GetUserByEmail(email)
	name := ""
	if user != nil {
		name = user.Name
	}

	return c.Render("dashboard", fiber.Map{
		"Title":      "Dashboard",
		"Email":      email,
		"Name":       name,
		"Workspaces": workspaces,
	})
}
