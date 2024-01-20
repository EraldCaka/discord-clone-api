package handlers

import (
	"errors"
	"github.com/EraldCaka/discord-clone-api/db"
	"github.com/EraldCaka/discord-clone-api/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChannelHandler struct {
	channelStore db.ChannelStore
}

func NewChannelHandler(channelStore db.ChannelStore) *ChannelHandler {
	return &ChannelHandler{
		channelStore: channelStore,
	}
}

func (h *ChannelHandler) HandleGetChannels(c *fiber.Ctx) error {
	channels, err := h.channelStore.GetChannels(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(channels)
}
func (h *ChannelHandler) HandleGetChannel(c *fiber.Ctx) error {
	var id = c.Params("id")
	channel, err := h.channelStore.GetChannelByID(c.Context(), id)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return c.JSON(map[string]string{"error": "not found"})
	}
	return c.JSON(channel)
}

func (h *ChannelHandler) HandleCreateChannel(c *fiber.Ctx) error {
	var params types.CreateChannelParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if validate := params.Validate(); len(validate) > 0 {
		return c.JSON(validate)
	}
	channel, err := types.NewChannel(params)
	if err != nil {
		return err
	}
	createdChannel, err := h.channelStore.CreateChannel(c.Context(), channel)
	return c.JSON(createdChannel)
}
func (h *ChannelHandler) HandleDeleteChannel(c *fiber.Ctx) error {
	client, _ := mongo.Connect(c.Context(), options.Client().ApplyURI(db.MONGODB))
	channelID := c.Params("id")

	if err := h.channelStore.DeleteChannel(c.Context(), channelID, client); err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": channelID})
}
