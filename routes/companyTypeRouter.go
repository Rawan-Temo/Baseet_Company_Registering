package routes

import (
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupCompanyTypeRoutes(api fiber.Router) {
	typeRoute := api.Group("/types")
	typeRoute.Get("/", handlers.AllCompanyTypes)
	typeRoute.Post("/", handlers.CreateCompanyType)
	typeRoute.Get(":id", handlers.GetCompanyType)
	typeRoute.Patch(":id", handlers.UpdateCompanyType)
	typeRoute.Delete(":id", handlers.DeleteCompanyType)
}
