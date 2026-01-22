package company_routes

import (
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	handlers "github.com/Rawan-Temo/Baseet_Company_Registering.git/handlers/company"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/middlewares"
	company_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/company"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"github.com/gofiber/fiber/v2"
)

func SetupTradingActivityRoutes(api fiber.Router) {

	tradingActivity := api.Group("/tradingActivities")

	// Protected routes
	tradingActivity.Use(middlewares.IsAuthenticated)
	tradingActivity.Patch("/delete-many", utils.DeleteMany(database.DB, company_models.TradingActivity{}))

	tradingActivity.Get("/", handlers.AllTradingActivity)
	tradingActivity.Post("/", handlers.CreateTradingActivity)
	tradingActivity.Get("/:id", handlers.GetTradingActivityByID)
	tradingActivity.Patch("/:id", handlers.UpdateTradingActivity)
	tradingActivity.Delete("/:id", handlers.DeleteTradingActivity)

}
