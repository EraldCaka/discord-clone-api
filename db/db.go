package db

import (
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
)

var (
	NAME    = os.Getenv("MONGO_DB_NAME")
	MONGODB = os.Getenv("MONGO_DB_URL")
)

type Store struct {
	User    UserStore
	Server  ServerStore
	Channel ChannelStore
	Message MessageStore
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	NAME = os.Getenv("MONGO_DB_NAME")
	MONGODB = os.Getenv("MONGO_DB_URL")
}

type MongoStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoStore(client *mongo.Client, collection string) *MongoStore {
	return &MongoStore{
		client: client,
		coll:   client.Database(NAME).Collection(collection),
	}
}
