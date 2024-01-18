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
		mongoDBName   = db.DBNAME
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
	PopulateUserTable(store)
}

func PopulateUserTable(store *db.Store) {
	for i := 0; i < 20; i++ {
		user := fixtures.AddUser(store, "user"+strconv.Itoa(i), "test1234"+strconv.Itoa(i), "I am a user", "test"+strconv.Itoa(i)+"@gmail.com")
		log.Println(user.Username, "was created successfully")
		PopulateServerTable(store, user.ID)
	}

}
func PopulateServerTable(store *db.Store, userID primitive.ObjectID) {
	server := fixtures.AddServer(store, "server", userID, "europe", "afk-channel", "europe")
	//TODO properly connect server seeding with the user seeding
	log.Println(server.ServerName, "was created successfully")
}
