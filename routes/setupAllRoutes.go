package routes

import (
	auth_routes "github.com/Rawan-Temo/Baseet_Company_Registering.git/routes/auth"
	company_routes "github.com/Rawan-Temo/Baseet_Company_Registering.git/routes/company"
	general_routes "github.com/Rawan-Temo/Baseet_Company_Registering.git/routes/general"
	media_routes "github.com/Rawan-Temo/Baseet_Company_Registering.git/routes/media"
	"github.com/gofiber/fiber/v2"
)

func SetupAllRoutes(app *fiber.App) {
	// Register routes (pass db pointer)
	app.Route("api/v1", auth_routes.SetupUserRoutes)
	app.Route("api/v1", auth_routes.SetupLicenseRoutes)
	app.Route("api/v1", company_routes.SetupCompanyActivityRoutes)
	app.Route("api/v1", company_routes.SetupTradingActivityRoutes)
	app.Route("api/v1", company_routes.SetupPeopleRoutes)
	app.Route("api/v1", company_routes.SetupCompanyRoutes)
	app.Route("api/v1", general_routes.SetupCompanyTypeRoutes)
	app.Route("api/v1", general_routes.SetupOfficeRoutes)
	app.Route("api/v1", media_routes.SetupMediaTypeRoutes)

}
