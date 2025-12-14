package company_routes

import (
	handlers "github.com/Rawan-Temo/Baseet_Company_Registering.git/handlers/company"
	"github.com/gofiber/fiber/v2"
)

func SetupCompanyRoutes(api fiber.Router) {
	company := api.Group("/companies")
	// Public routes
	company.Post("/", handlers.CreateCompany)
	// Protected routes
	// company.Use(middlewares.IsAuthenticated)

	company.Get("/", handlers.AllCompanies)
	company.Get("/:id", handlers.SingleCompany)
	company.Patch("/:id", handlers.UpdateCompany)
	company.Delete("/:id", handlers.DeleteCompany)

}
