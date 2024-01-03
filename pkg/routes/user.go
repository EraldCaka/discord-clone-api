package routes

import (
	"github.com/EraldCaka/discord-clone-api/handlers"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(a *fiber.App, userHandler *handlers.UserHandler) {
	var route = a.Group("/discord/api/v1")

	route.Get("/user/:id", userHandler.HandleGetUser)
	route.Put("/user/:id", userHandler.HandlePutUser)
	route.Delete("/user/:id", userHandler.HandleDeleteUser)
	route.Post("/user", userHandler.HandleCreateUser)
	route.Get("/user", userHandler.HandleGetUsers)
}
