package routes

import (
	"github.com/EraldCaka/discord-clone-api/handlers"
	"github.com/gofiber/fiber/v2"
)

func ChannelRoutes(a *fiber.App, channelHandler *handlers.ChannelHandler, route fiber.Router) {
	route.Get("/channel/:id", channelHandler.HandleGetChannel)
	route.Delete("/channel/:id", channelHandler.HandleDeleteChannel)
	route.Post("/channel", channelHandler.HandleCreateChannel)
	route.Get("/channel", channelHandler.HandleGetChannels)
}
