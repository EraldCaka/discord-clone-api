package types

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const messageLen = 250

type Message struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Content         string             `bson:"content" json:"content"`
	Author          User               `bson:"author" json:"author"`
	Channel         Channel            `bson:"channel" json:"channel"`
	Attachments     []Attachment       `bson:"attachments" json:"attachments"`
	Reactions       []Reaction         `bson:"reactions" json:"reactions"`
	Mentions        []User             `bson:"mentions" json:"mentions"`
	MentionRoles    []Role             `bson:"mentionRoles" json:"mentionRoles"`
	MentionEveryone bool               `bson:"mentionEveryone" json:"mentionEveryone"`
	Pinned          bool               `bson:"pinned" json:"pinned"`
}

func (params CreateMessageParams) Validate() map[string]string {
	//TODO include all errors for each use case like Ownership, Members, Channels, Roles and AFKChannel
	errors := map[string]string{}
	if len(params.Content) > messageLen {
		errors["description"] = fmt.Sprintf("description length should be at least %d characters", messageLen)
	}
	return errors
}

type CreateMessageParams struct {
	Content         string       `json:"content"`
	Author          User         `json:"author"`
	Channel         Channel      `json:"channel"`
	Attachments     []Attachment `json:"attachments"`
	Reactions       []Reaction   `json:"reactions"`
	Mentions        []User       `json:"mentions"`
	MentionRoles    []Role       `json:"mentionRoles"`
	MentionEveryone bool         `json:"mentionEveryone"`
	Pinned          bool         `json:"pinned"`
}

func NewMessage(params CreateMessageParams) (*Message, error) {
	return &Message{
		Content:         params.Content,
		Author:          params.Author,
		Channel:         params.Channel,
		Attachments:     params.Attachments,
		Reactions:       params.Reactions,
		Mentions:        params.Mentions,
		MentionRoles:    params.MentionRoles,
		MentionEveryone: params.MentionEveryone,
		Pinned:          params.Pinned,
	}, nil
}
