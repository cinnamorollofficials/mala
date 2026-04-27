package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hadigunawan/mala/internal/routes"
	"github.com/hadigunawan/mala/pkg/database"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Connect to database
	database.ConnectDB()
	defer database.DB.Close()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Mala LLM Gateway",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// Setup routes
	routes.SetupRoutes(app)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
