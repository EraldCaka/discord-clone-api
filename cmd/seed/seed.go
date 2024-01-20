package main

import (
	"context"
	"github.com/EraldCaka/discord-clone-api/db"
	"github.com/EraldCaka/discord-clone-api/db/fixtures"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
)

func main() {
	var (
		ctx           = context.Background()
		mongoEndpoint = db.MONGODB
		mongoDBName   = db.NAME
	)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(mongoDBName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	userStore := db.NewMongoUserStore(client)
	serverStore := db.NewMongoServerStore(client, userStore)
	channelStore := db.NewMongoChannelStore(client, serverStore)
	messageStore := db.NewMongoMessageStore(client, channelStore, userStore)
	store := &db.Store{
		User:    userStore,
		Server:  serverStore,
		Channel: channelStore,
		Message: messageStore,
	}
	PopulateTables(store)

}
func PopulateTables(store *db.Store) {

	for i := 0; i < 5; i++ {
		userID := PopulateUserTable(store, i)
		serverID := PopulateServerTable(store, userID)
		for x := 0; x < 3; x++ {
			channelID := PopulateChannelTable(store, serverID)
			for j := 0; j < 5; j++ {
				PopulateMessageTable(store, channelID, userID)
			}
		}
	}
}
func PopulateUserTable(store *db.Store, i int) primitive.ObjectID {
	user := fixtures.AddUser(store, "user"+strconv.Itoa(i), "test1234"+strconv.Itoa(i), "I am a user", "test"+strconv.Itoa(i)+"@gmail.com")

	log.Println(user.Username, "was created successfully!")
	return user.ID
}

func PopulateServerTable(store *db.Store, userID primitive.ObjectID) primitive.ObjectID {
	server := fixtures.AddServer(store, "server", userID, "europe", "afk-channel", "europe")
	log.Println(server.ServerName, "was created successfully!")
	return server.ID
}

func PopulateChannelTable(store *db.Store, serverID primitive.ObjectID) primitive.ObjectID {
	channel := fixtures.AddChannel(store, serverID, "general", true, "everyone can talk", false)
	log.Println(channel.ChannelName, "was created successfully!")
	return channel.ID
}

func PopulateMessageTable(store *db.Store, channelID primitive.ObjectID, userID primitive.ObjectID) primitive.ObjectID {
	message := fixtures.AddMessage(store, channelID, userID, "hi i am a user", false, false)
	log.Println(message.ID, "was created successfully")
	return message.ID
}
