package routes

import (
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/handlers"
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
