package auth_routes

import (
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	handlers "github.com/Rawan-Temo/Baseet_Company_Registering.git/handlers/auth"
	auth_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/auth"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"github.com/gofiber/fiber/v2"
)








func SetupLicenseRoutes(api fiber.Router) {
	license := api.Group("/licenses")
		// Protected routes
	// company.Use(middlewares.IsAuthenticated)
	license.Patch("/delete-many", utils.DeleteMany(database.DB, auth_models.License{}))

	license.Get("/" , handlers.GetAllLicenses)
	license.Post("/", handlers.CreateLicense)
	license.Get("/:id" , handlers.GetLicenseByID)
	license.Patch("/:id" , handlers.UpdateLicense)
	license.Delete("/:id" , handlers.DeleteLicense)


}