package routes

import (
	"github.com/EraldCaka/discord-clone-api/handlers"
	"github.com/gofiber/fiber/v2"
)

func MessageRoutes(a *fiber.App, messageHandler *handlers.MessageHandler) {
	var route = a.Group("/discord/api/v1")

	route.Get("/message/:id", messageHandler.HandleGetMessage)
	route.Delete("/message/:id", messageHandler.HandleDeleteMessage)
	route.Post("/message", messageHandler.HandleCreateMessage)
	route.Get("/message", messageHandler.HandleGetMessages)
}
