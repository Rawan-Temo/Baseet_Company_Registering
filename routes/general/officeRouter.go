package general_routes

import (
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	handlers "github.com/Rawan-Temo/Baseet_Company_Registering.git/handlers/general"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/middlewares"
	general_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/general"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"github.com/gofiber/fiber/v2"
)

func SetupOfficeRoutes(api fiber.Router) {
	officeRouter := api.Group("/offices")

	officeRouter.Get("/", handlers.AllOffices)
	officeRouter.Use(middlewares.IsAuthenticated)
	officeRouter.Patch("/delete-many", utils.DeleteMany(database.DB, general_models.Office{}))
	officeRouter.Post("/", handlers.CreateOffice)
	officeRouter.Get("/:id", handlers.GetOffice)
	officeRouter.Patch("/:id", handlers.UpdateOffice)
	officeRouter.Delete("/:id", handlers.DeleteOffice)
}
