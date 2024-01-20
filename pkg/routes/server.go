package routes

import (
	"github.com/EraldCaka/discord-clone-api/handlers"
	"github.com/gofiber/fiber/v2"
)

func ServerRoutes(app *fiber.App, serverHandler *handlers.ServerHandler, route fiber.Router) {
	route.Get("/server/:id", serverHandler.HandleGetServer)
	route.Delete("/server/:id", serverHandler.HandleDeleteServer)
	route.Post("/server", serverHandler.HandleCreateServer)
	route.Get("/server", serverHandler.HandleGetServers)
}
