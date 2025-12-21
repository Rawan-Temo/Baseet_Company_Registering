package company_routes

import (
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	handlers "github.com/Rawan-Temo/Baseet_Company_Registering.git/handlers/company"
	company_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/company"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"github.com/gofiber/fiber/v2"
)

func SetupCompanyRoutes(api fiber.Router) {
	company := api.Group("/companies")
	// Public routes
	company.Post("/", handlers.CreateCompany)
	// Protected routes
	// company.Use(middlewares.IsAuthenticated)
	company.Patch("/delete-many", utils.DeleteMany(database.DB, company_models.Company{}))
	company.Get("/", handlers.AllCompanies)
	company.Get("/:id", handlers.SingleCompany)
	company.Patch("/:id", handlers.UpdateCompany)
	company.Delete("/:id", handlers.DeleteCompany)

}
