package routes

import (
	"waweb-v2/handlers"
	"waweb-v2/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Public routes
	app.Get("/", handlers.Index)
	app.Get("/login", handlers.LoginPage)
	app.Post("/login", handlers.Login)
	app.Get("/register", handlers.RegisterPage)
	app.Post("/register", handlers.Register)
	app.Post("/logout", handlers.Logout)

	// Protected routes
	app.Get("/dashboard", middleware.AuthRequired(), handlers.Dashboard)

	// Workspace routes
	app.Get("/workspace/new", middleware.AuthRequired(), handlers.NewWorkspace)
	app.Post("/workspace/create", middleware.AuthRequired(), handlers.CreateWorkspace)
	app.Get("/workspace/:id", middleware.AuthRequired(), handlers.WorkspaceDetail)
	app.Delete("/workspace/:id", middleware.AuthRequired(), handlers.DeleteWorkspace)

	// Template picker
	app.Get("/templates", middleware.AuthRequired(), handlers.ListTemplates)
	app.Get("/templates/:id", middleware.AuthRequired(), handlers.GetTemplate)

	// Builder
	app.Get("/builder/:id", middleware.AuthRequired(), handlers.Builder)
	app.Post("/builder/:id/save", middleware.AuthRequired(), handlers.SaveBuilder)
	app.Post("/builder/:id/publish", middleware.AuthRequired(), handlers.PublishSite)
	app.Post("/builder/:id/upload", middleware.AuthRequired(), handlers.UploadAsset)

	// Preview (live preview before publishing)
	app.Get("/preview/:id/*", handlers.PreviewSite)

	// Serve uploaded assets (images etc)
	app.Static("/uploads", "./uploads")

	// Published sites - served from host folder
	app.Get("/sites/:subdomain/*", handlers.ServeSite)
}
