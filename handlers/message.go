package handlers

import (
	"errors"
	"github.com/EraldCaka/discord-clone-api/db"
	"github.com/EraldCaka/discord-clone-api/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageHandler struct {
	messageStore db.MessageStore
}

func NewMessageHandler(messageStore db.MessageStore) *MessageHandler {
	return &MessageHandler{
		messageStore: messageStore,
	}
}

func (h *MessageHandler) HandleGetMessages(c *fiber.Ctx) error {
	messages, err := h.messageStore.GetMessages(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(messages)
}
func (h *MessageHandler) HandleGetMessage(c *fiber.Ctx) error {
	var id = c.Params("id")
	message, err := h.messageStore.GetMessageByID(c.Context(), id)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return c.JSON(map[string]string{"error": "not found"})
	}
	return c.JSON(message)
}

func (h *MessageHandler) HandleCreateMessage(c *fiber.Ctx) error {
	var params types.CreateMessageParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if validate := params.Validate(); len(validate) > 0 {
		return c.JSON(validate)
	}
	message, err := types.NewMessage(params)
	if err != nil {
		return err
	}
	createdMessage, err := h.messageStore.CreateMessage(c.Context(), message)
	return c.JSON(createdMessage)
}
func (h *MessageHandler) HandleDeleteMessage(c *fiber.Ctx) error {
	messageID := c.Params("id")
	if err := h.messageStore.DeleteMessage(c.Context(), messageID); err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": messageID})
}
