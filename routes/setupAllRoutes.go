package routes

import "github.com/gofiber/fiber/v2"

func SetupAllRoutes(app *fiber.App) {
	// Register routes (pass db pointer)
	app.Route("api/v1", SetupUserRoutes)
	app.Route("api/v1", SetupCompanyRoutes)
	app.Route("api/v1", SetupCompanyTypeRoutes)
	app.Route("api/v1", SetupOfficeRoutes)
}
