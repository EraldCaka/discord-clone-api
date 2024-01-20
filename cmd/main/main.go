package main

import (
	"context"
	"flag"
	"github.com/EraldCaka/discord-clone-api/db"
	"github.com/EraldCaka/discord-clone-api/handlers"
	"github.com/EraldCaka/discord-clone-api/pkg/errors"
	"github.com/EraldCaka/discord-clone-api/pkg/routes"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var config = fiber.Config{
	ErrorHandler: errors.ErrorHandler,
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.MONGODB))
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New(config)

	var (
		userStore    = db.NewMongoUserStore(client)
		serverStore  = db.NewMongoServerStore(client, userStore)
		channelStore = db.NewMongoChannelStore(client, serverStore)
		messageStore = db.NewMongoMessageStore(client, channelStore, userStore)
		store        = &db.Store{
			User:   userStore,
			Server: serverStore,
		}
		userHandler    = handlers.NewUserHandler(store)
		serverHandler  = handlers.NewServerHandler(serverStore)
		channelHandler = handlers.NewChannelHandler(channelStore)
		messageHandler = handlers.NewMessageHandler(messageStore)
	)
	routes.UserRoutes(app, userHandler)
	routes.ServerRoutes(app, serverHandler)
	routes.ChannelRoutes(app, channelHandler)
	routes.MessageRoutes(app, messageHandler)

	listenErr := app.Listen(*listenAddr)
	if listenErr != nil {
		return
	}
}

// docker run --name mongodb -d -p 27017:27017 mongo:latest
