package general_routes

import (
	handlers "github.com/Rawan-Temo/Baseet_Company_Registering.git/handlers/general"
	"github.com/gofiber/fiber/v2"
)

func SetupOfficeRoutes(api fiber.Router) {
	officeRoute := api.Group("/offices")
	officeRoute.Get("/", handlers.AllOffices)
	officeRoute.Post("/", handlers.CreateOffice)
	officeRoute.Get(":id", handlers.GetOffice)
	officeRoute.Patch(":id", handlers.UpdateOffice)
	officeRoute.Delete(":id", handlers.DeleteOffice)
}
