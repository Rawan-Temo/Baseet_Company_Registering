package company_routes

import (
	handlers "github.com/Rawan-Temo/Baseet_Company_Registering.git/handlers/company"
	"github.com/gofiber/fiber/v2"
)




func SetupPeopleRoutes(api fiber.Router){
	people := api.Group("/people")
		// Protected routes
	// company.Use(middlewares.IsAuthenticated)

	people.Get("/" , handlers.GetAllPeople)
	people.Post("/", handlers.CreatePerson)
	people.Get("/:id" , handlers.GetPersonByID	)
	people.Patch("/:id" , handlers.UpdatePerson)
	people.Delete("/:id" , handlers.DeletePerson)

	
}