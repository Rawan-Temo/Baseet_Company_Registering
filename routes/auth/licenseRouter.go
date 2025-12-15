package auth_routes

import (
	handlers "github.com/Rawan-Temo/Baseet_Company_Registering.git/handlers/auth"
	"github.com/gofiber/fiber/v2"
)








func SetupLicenseRoutes(api fiber.Router) {
	license := api.Group("/licenses")
		// Protected routes
	// company.Use(middlewares.IsAuthenticated)

	license.Get("/" , handlers.GetAllLicenses)
	license.Post("/", handlers.CreateLicense)
	license.Get("/:id" , handlers.GetLicenseByID)
	license.Patch("/:id" , handlers.UpdateLicense)
	license.Delete("/:id" , handlers.DeleteLicense)


}