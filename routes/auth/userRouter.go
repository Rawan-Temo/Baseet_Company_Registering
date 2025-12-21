package auth_routes

import (
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	handlers "github.com/Rawan-Temo/Baseet_Company_Registering.git/handlers/auth"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/middlewares"
	auth_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/auth"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(api fiber.Router) {
	database:=database.DB
	user := api.Group("/users")
	// Public routes
	user.Post("/login", handlers.Login)
	user.Patch("/delete-many" ,utils.DeleteMany(database  , auth_models.User{})) 
	// Protected routes
	user.Post("/", handlers.CreateUser)
	user.Use(middlewares.IsAuthenticated)
	user.Get("/profile",  handlers.GetUserFromToken)
	
	user.Get("/",  handlers.AllUsers)
	user.Get("/:id", handlers.SingleUser)
	user.Patch("/:id", handlers.UpdateUser)
	user.Delete("/:id", handlers.DeleteUser)

}
