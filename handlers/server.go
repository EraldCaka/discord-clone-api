package handlers

import (
	"errors"
	"github.com/EraldCaka/discord-clone-api/db"
	"github.com/EraldCaka/discord-clone-api/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServerHandler struct {
	serverStore db.ServerStore
}

func NewServerHandler(serverStore db.ServerStore) *ServerHandler {
	return &ServerHandler{
		serverStore: serverStore,
	}
}

func (h *ServerHandler) HandleGetServers(c *fiber.Ctx) error {
	servers, err := h.serverStore.GetServers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(servers)
}

func (h *ServerHandler) HandleGetServer(c *fiber.Ctx) error {
	var id = c.Params("id")
	server, err := h.serverStore.GetServerByID(c.Context(), id)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return c.JSON(map[string]string{"error": "not found"})
	}
	return c.JSON(server)
}

func (h *ServerHandler) HandleCreateServer(c *fiber.Ctx) error {
	var params types.CreateServerParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if validate := params.Validate(); len(validate) > 0 {
		return c.JSON(validate)
	}
	server, err := types.NewServer(params)
	if err != nil {
		return err
	}
	createdServer, err := h.serverStore.CreateServer(c.Context(), server)
	return c.JSON(createdServer)
}

func (h *ServerHandler) HandleDeleteServer(c *fiber.Ctx) error {
	serverID := c.Params("id")
	if err := h.serverStore.DeleteServer(c.Context(), serverID); err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": serverID})
}
