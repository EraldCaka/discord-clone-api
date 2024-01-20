package fixtures

import (
	"context"
	"github.com/EraldCaka/discord-clone-api/db"
	"github.com/EraldCaka/discord-clone-api/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

func AddUser(store *db.Store, userName, password, description, email string) *types.User {
	user, err := types.NewUser(types.CreateUserParams{
		Username:    userName,
		Password:    password,
		Description: description,
		Email:       email, OwnedServers: make([]primitive.ObjectID, 0),
	})

	insertedUser, err := store.User.CreateUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}

func AddServer(store *db.Store, serverName string, userID primitive.ObjectID, region, afkChannel, description string) *types.Server {
	server, _ := types.NewServer(types.CreateServerParams{
		ServerName:  serverName,
		UserID:      userID,
		Members:     []types.User{},
		Roles:       []types.Role{},
		Channels:    make([]primitive.ObjectID, 0),
		Region:      region,
		AfkChannel:  afkChannel,
		Description: description,
	})
	insertedServer, err := store.Server.CreateServer(context.TODO(), server)
	if err != nil {
		log.Fatal(err)
	}
	return insertedServer
}

func AddChannel(store *db.Store, serverID primitive.ObjectID, channelName string, chanType bool, description string, nsfw bool) *types.Channel {
	channel, _ := types.NewChannel(types.CreateChannelParams{
		ServerID:    serverID,
		Messages:    make([]primitive.ObjectID, 0),
		ChannelName: channelName,
		Type:        chanType,
		Description: description,
		Nsfw:        nsfw,
	})
	insertedChannel, err := store.Channel.CreateChannel(context.TODO(), channel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedChannel
}
func AddMessage(store *db.Store, channelID primitive.ObjectID, userID primitive.ObjectID, content string, mentionEveryone, pinned bool) *types.Message {
	message, _ := types.NewMessage(types.CreateMessageParams{
		ChannelID:       channelID,
		UserID:          userID,
		Content:         content,
		MentionEveryone: mentionEveryone,
		Pinned:          pinned,
	})
	insertedMessage, err := store.Message.CreateMessage(context.TODO(), message)
	if err != nil {
		log.Fatal(err)
	}
	return insertedMessage
}
