package routes

import (
	"github.com/EraldCaka/discord-clone-api/handlers"
	"github.com/gofiber/fiber/v2"
)

func ServerRoutes(a *fiber.App, serverHandler *handlers.ServerHandler) {
	var route = a.Group("/discord/api/v1")

	route.Get("/server/:id", serverHandler.HandleGetServer)
	route.Delete("/server/:id", serverHandler.HandleDeleteServer)
	route.Post("/server", serverHandler.HandleCreateServer)
	route.Get("/server", serverHandler.HandleGetServers)
}
