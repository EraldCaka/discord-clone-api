package types

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const descriptionLimit = 100

type CreateServerParams struct {
	ServerName  string             `json:"serverName" validate:"required"`
	UserID      primitive.ObjectID `json:"userID,omitempty"`
	Members     []User             `json:"members"`
	Channels    []Channel          `json:"channels"`
	Roles       []Role             `json:"roles"`
	Region      string             `json:"region" validate:"required"`
	AfkChannel  Channel            `json:"afkChannel"`
	Description string             `json:"description"`
}

func (params CreateServerParams) Validate() map[string]string {
	//TODO include all errors for each use case like Ownership, Members, Channels, Roles and AFKChannel
	errors := map[string]string{}
	if len(params.Description) > descriptionLen {
		errors["description"] = fmt.Sprintf("description length should be at least %d characters", descriptionLen)
	}
	return errors
}

type Server struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	ServerName  string             `bson:"serverName" json:"serverName"`
	Members     []User             `bson:"members" json:"members"`
	Channels    []Channel          `bson:"channels" json:"channels"`
	Roles       []Role             `bson:"roles" json:"roles"`
	Region      string             `bson:"regions" json:"regions"`
	AfkChannel  Channel            `bson:"afkChannel" json:"afkChannel"`
	Description string             `bson:"description" json:"description"`
}

func NewServer(params CreateServerParams) (*Server, error) {
	return &Server{
		ServerName:  params.ServerName,
		UserID:      params.UserID,
		Members:     params.Members,
		Channels:    params.Channels,
		Roles:       params.Roles,
		Region:      params.Region,
		AfkChannel:  params.AfkChannel,
		Description: params.Description,
	}, nil
}
