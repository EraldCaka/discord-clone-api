package types

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ChannelID primitive.ObjectID `bson:"channelID,omitempty" json:"channelID,omitempty"`
	UserID    primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	Content   string             `bson:"content" json:"content"`
	//Channel         Channel            `bson:"channel" json:"channel"`
	//Attachments     []Attachment `bson:"attachments" json:"attachments"`
	//Reactions       []Reaction   `bson:"reactions" json:"reactions"`
	//Mentions        []User       `bson:"mentions" json:"mentions"`
	//MentionRoles    []Role       `bson:"mentionRoles" json:"mentionRoles"`
	MentionEveryone bool `bson:"mentionEveryone" json:"mentionEveryone"`
	Pinned          bool `bson:"pinned" json:"pinned"`
}

func (params CreateMessageParams) Validate() map[string]string {
	//TODO
	errors := map[string]string{}
	if len(params.Content) > messageLen {
		errors["content"] = fmt.Sprintf("content length should be at least %d characters", messageLen)
	}
	return errors
}

type CreateMessageParams struct {
	ChannelID primitive.ObjectID `json:"channelID"`
	UserID    primitive.ObjectID `json:"userID"`
	Content   string             `json:"content"`

	//Channel         Channel            `json:"channel"`
	//Attachments     []Attachment `json:"attachments"`
	//Reactions       []Reaction   `json:"reactions"`
	//Mentions        []User       `json:"mentions"`
	//MentionRoles    []Role       `json:"mentionRoles"`
	MentionEveryone bool `json:"mentionEveryone"`
	Pinned          bool `json:"pinned"`
}

func NewMessage(params CreateMessageParams) (*Message, error) {
	return &Message{
		ChannelID: params.ChannelID,
		UserID:    params.UserID,
		Content:   params.Content,
		//Channel:         params.Channel,
		//Attachments:     params.Attachments,
		//Reactions:       params.Reactions,
		//Mentions:        params.Mentions,
		//MentionRoles:    params.MentionRoles,
		MentionEveryone: params.MentionEveryone,
		Pinned:          params.Pinned,
	}, nil
}
