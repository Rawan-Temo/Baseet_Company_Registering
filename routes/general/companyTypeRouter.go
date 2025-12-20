package general_routes

import (
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	handlers "github.com/Rawan-Temo/Baseet_Company_Registering.git/handlers/general"
	general_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/general"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"

	"github.com/gofiber/fiber/v2"
)

func SetupCompanyTypeRoutes(api fiber.Router) {
	// var controllers  = utils.CreateControllers(models.CompanyType{})

	typeRouter := api.Group("/types")
	typeRouter.Patch("/delete-many", utils.DeleteMany(database.DB , general_models.CompanyType{}))

	typeRouter.Get("/", handlers.AllCompanyTypes)
	typeRouter.Post("/", handlers.CreateCompanyType)
	typeRouter.Get("/:id", handlers.GetCompanyType)
	typeRouter.Patch("/:id", handlers.UpdateCompanyType)
	typeRouter.Delete("/:id", handlers.DeleteCompanyType)
}
