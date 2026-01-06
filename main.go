package main

import (
	"log"
	"os"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

// TODOO : company activities , trading activities , people routes protect them with authentication middleware
// TODOO : handlers for licenses  companyActivities middle table check media handlers maybe u screwed up something
// TODOO : check compnay category functionality
// and create the regiester end point
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using defaults")
	}
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	app.Static("/", "./uploads")
	// Connect DB
	database.ConnectDB()
	routes.SetupAllRoutes(app)



	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("🚀 Server running on http://localhost:%s", port)
	app.Listen(":" + port)
}
