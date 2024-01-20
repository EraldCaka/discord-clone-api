package routes

import (
	"github.com/EraldCaka/discord-clone-api/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, authHandler *middleware.AuthHandler, route fiber.Router) {
	route.Post("/auth", authHandler.HandleAuthenticate)

}
