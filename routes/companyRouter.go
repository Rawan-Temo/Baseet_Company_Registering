package routes

import (
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupCompanyRoutes(api fiber.Router) {
	company := api.Group("/companies")
	// Public routes
	// Protected routes
	// company.Use(middlewares.IsAuthenticated)
	
	company.Post("/", handlers.CreateCompany)
	company.Get("/",  handlers.AllCompanies)
	company.Get("/:id", handlers.SingleCompany)
	company.Patch("/:id", handlers.UpdateCompany)
	company.Delete("/:id", handlers.DeleteCompany)

}
