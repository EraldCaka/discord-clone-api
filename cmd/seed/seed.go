package main

import (
	"context"
	"github.com/EraldCaka/discord-clone-api/db"
	"github.com/EraldCaka/discord-clone-api/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
)

var ()

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
	store := &db.Store{
		User:    userStore,
		Server:  serverStore,
		Channel: channelStore,
		Message: db.NewMongoMessageStore(client, channelStore, userStore),
	}
	PopulateUserTable(store)

}

func PopulateUserTable(store *db.Store) {
	for i := 0; i < 20; i++ {
		user := fixtures.AddUser(store, "user"+strconv.Itoa(i), "test1234"+strconv.Itoa(i), "I am a user", "test"+strconv.Itoa(i)+"@gmail.com")
		log.Println(user.Username, "was created successfully")
	}

}
