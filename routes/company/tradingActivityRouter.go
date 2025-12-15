package company_routes

import (
	handlers "github.com/Rawan-Temo/Baseet_Company_Registering.git/handlers/company"
	"github.com/gofiber/fiber/v2"
)





func SetupTradingActivityRoutes(api fiber.Router){
	
	tradingActivity := api.Group("/tradingActivities")

		// Protected routes
	// company.Use(middlewares.IsAuthenticated)

	tradingActivity.Get("/" , handlers.AllTradingActivity)
	tradingActivity.Post("/", handlers.CreateTradingActivity)
	tradingActivity.Get("/:id" , handlers.GetTradingActivityByID	)
	tradingActivity.Patch("/:id" , handlers.UpdateTradingActivity)
	tradingActivity.Delete("/:id" , handlers.DeleteTradingActivity)
	
}