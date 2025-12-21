package company_routes

import (
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	handlers "github.com/Rawan-Temo/Baseet_Company_Registering.git/handlers/company"
	company_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/company"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"github.com/gofiber/fiber/v2"
)




func SetupCompanyActivityRoutes(api fiber.Router){
	companyActivity := api.Group("/companyActivities")
		// Protected routes
	// company.Use(middlewares.IsAuthenticated)
    companyActivity.Patch("/delete-many" ,utils.DeleteMany(database.DB  , company_models.CompanyActivity{})) 
	companyActivity.Get("/" , handlers.GetAllCompanyActivities)
	companyActivity.Post("/", handlers.CreateCompanyActivity)
	companyActivity.Get("/:id" , handlers.GetCompanyActivityByID)
	companyActivity.Patch("/:id" , handlers.UpdateCompanyActivity)
	companyActivity.Delete("/:id" , handlers.DeleteCompanyActivity)

	
}