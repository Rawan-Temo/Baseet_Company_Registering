package auth_routes

import (
	handlers "github.com/Rawan-Temo/Baseet_Company_Registering.git/handlers/auth"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(api fiber.Router) {
	user := api.Group("/users")
	// Public routes
	user.Post("/login", handlers.Login)
	// Protected routes
	user.Post("/", handlers.CreateUser)
	// user.Use(middlewares.IsAuthenticated)
	
	user.Get("/",  handlers.AllUsers)
	user.Get("/:id", handlers.SingleUser)
	user.Patch("/:id", handlers.UpdateUser)
	user.Delete("/:id", handlers.DeleteUser)

}
