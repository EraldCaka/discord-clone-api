package routes

import (
	"github.com/EraldCaka/discord-clone-api/handlers"
	"github.com/gofiber/fiber/v2"
)

func ChannelRoutes(a *fiber.App, channelHandler *handlers.ChannelHandler) {
	var route = a.Group("/discord/api/v1")

	route.Get("/channel/:id", channelHandler.HandleGetChannel)
	route.Delete("/channel/:id", channelHandler.HandleDeleteChannel)
	route.Post("/channel", channelHandler.HandleCreateChannel)
	route.Get("/channel", channelHandler.HandleGetChannels)
}
