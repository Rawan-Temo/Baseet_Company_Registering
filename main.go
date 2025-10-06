package main

import (
	"log"
	"os"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using defaults")
	}

	app := fiber.New()
	app.Use(logger.New())

	// Connect DB

	database.ConnectDB()
	// Register routes (pass db pointer)
	app.Route("api/v1" , routes.SetupUserRoutes)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("🚀 Server running on http://localhost:%s", port)
	app.Listen(":" + port)
}
