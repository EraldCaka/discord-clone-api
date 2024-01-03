package types

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateChannelParams struct {
	ChannelName string `json:"channelName"`
	Description string `json:"description"`
	Nsfw        bool   `json:"nsfw"`
}

func (p UpdateChannelParams) ToBSON() bson.M {
	m := bson.M{}
	if len(p.ChannelName) > 0 {
		m["channelName"] = p.ChannelName
	}
	if len(p.Description) > 0 {
		m["description"] = p.Description
	}
	if p.Nsfw {
		m["nsfw"] = true
	} else if p.Nsfw == false {
		m["nsfw"] = false
	}
	return m
}

type CreateChannelParams struct {
	ServerID    primitive.ObjectID   `json:"serverID"`
	Users       []primitive.ObjectID `json:"users"`
	ChannelName string               `json:"channelName"`
	Type        bool                 `json:"type"`
	Description string               `json:"description"`
	Nsfw        bool                 `json:"nsfw"`
}
type Channel struct {
	ID          primitive.ObjectID   `bson:"_if,omitempty" json:"id,omitempty"`
	ServerID    primitive.ObjectID   `bson:"serverID,omitempty" json:"serverID,omitempty"`
	Users       []primitive.ObjectID `bson:"users" json:"users"`
	ChannelName string               `bson:"channelName" json:"channelName"`
	Type        bool                 `bson:"type" json:"type"`
	Description string               `bson:"description" json:"description"`
	Nsfw        bool                 `bson:"nsfw" json:"nsfw"`
}

func (params CreateChannelParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.Description) > descriptionLen {
		errors["description"] = fmt.Sprintf("Channel description length should be at least %d characters", descriptionLen)
	}
	return errors
}
func NewChannel(params CreateChannelParams) (*Channel, error) {
	return &Channel{
		ServerID:    params.ServerID,
		Users:       params.Users,
		ChannelName: params.ChannelName,
		Type:        params.Type,
		Description: params.Description,
		Nsfw:        params.Nsfw,
	}, nil
}

/*
name: str
type: bool
description: str
position: int
category: Channel
nsfw: bool

*/
