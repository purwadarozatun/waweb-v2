package main

import (
	"encoding/json"
	"html/template"
	"log"
	"waweb-v2/config"
	"waweb-v2/database"
	"waweb-v2/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	htmlEngine "github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load config
	cfg := config.LoadConfig()

	// Initialize database
	database.Init()

	// Create templates engine
	engine := htmlEngine.New("./views", ".html")
	engine.Reload(true)
	engine.Debug(true)

	// Add template functions
	engine.AddFunc("toJSON", func(v interface{}) template.JS {
		bytes, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return template.JS("{}")
		}
		return template.JS(bytes)
	})

	// Create Fiber app
	app := fiber.New(fiber.Config{
		Views:     engine,
		BodyLimit: 10 * 1024 * 1024, // 10 MB — allows up to 5 MB image uploads
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// Static files
	app.Static("/static", "./static")
	app.Static("/assets", "./assets")

	// Setup routes
	routes.SetupRoutes(app)

	// Start server
	port := cfg.Port
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
