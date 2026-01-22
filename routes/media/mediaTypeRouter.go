package media_routes

import (
	handlers "github.com/Rawan-Temo/Baseet_Company_Registering.git/handlers/media"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupMediaTypeRoutes(router fiber.Router) {
	mediaType := router.Group("/media-types")

	mediaType.Use(middlewares.IsAuthenticated)
	mediaType.Get("/", handlers.GetAllMediaTypes)
	mediaType.Post("/", handlers.CreateMediaType)
	mediaType.Get("/:id", handlers.GetMediaTypeByID)
	mediaType.Put("/:id", handlers.UpdateMediaType)
	mediaType.Delete("/:id", handlers.DeleteMediaType)
}
