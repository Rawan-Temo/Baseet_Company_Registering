package company_routes

import (
	handlers "github.com/Rawan-Temo/Baseet_Company_Registering.git/handlers/company"
	"github.com/gofiber/fiber/v2"
)




func SetupCompanyActivityRoutes(api fiber.Router){
	companyActivity := api.Group("/companyActivities")
		// Protected routes
	// company.Use(middlewares.IsAuthenticated)

	companyActivity.Get("/" , handlers.GetAllCompanyActivities)
	companyActivity.Post("/", handlers.CreateCompanyActivity)
	companyActivity.Get("/:id" , handlers.GetCompanyActivityByID)
	companyActivity.Patch("/:id" , handlers.UpdateCompanyActivity)
	companyActivity.Delete("/:id" , handlers.DeleteCompanyActivity)

	
}