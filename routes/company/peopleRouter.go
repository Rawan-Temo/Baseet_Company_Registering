package company_routes

import (
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/database"
	handlers "github.com/Rawan-Temo/Baseet_Company_Registering.git/handlers/company"
	company_models "github.com/Rawan-Temo/Baseet_Company_Registering.git/models/company"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"github.com/gofiber/fiber/v2"
)




func SetupPeopleRoutes(api fiber.Router){
	people := api.Group("/people")
		// Protected routes
	// company.Use(middlewares.IsAuthenticated)
	people.Patch("/delete-many", utils.DeleteMany(database.DB, company_models.People{}))

	people.Get("/" , handlers.GetAllPeople)
	people.Post("/", handlers.CreatePerson)
	people.Get("/:id" , handlers.GetPersonByID	)
	people.Patch("/:id" , handlers.UpdatePerson)
	people.Delete("/:id" , handlers.DeletePerson)

	
}