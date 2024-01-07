package main

import (
	"context"
	"flag"
	"github.com/EraldCaka/discord-clone-api/db"
	"github.com/EraldCaka/discord-clone-api/handlers"
	"github.com/EraldCaka/discord-clone-api/pkg/routes"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var configs = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.MONGODB))
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New(configs)

	var (
		userStore    = db.NewMongoUserStore(client)
		serverStore  = db.NewMongoServerStore(client, userStore)
		channelStore = db.NewMongoChannelStore(client, serverStore)
		store        = &db.Store{
			User:   userStore,
			Server: serverStore,
		}
		userHandler    = handlers.NewUserHandler(store)
		serverHandler  = handlers.NewServerHandler(serverStore)
		channelHandler = handlers.NewChannelHandler(channelStore)
	)
	routes.UserRoutes(app, userHandler)
	routes.ServerRoutes(app, serverHandler)
	routes.ChannelRoutes(app, channelHandler)

	listenErr := app.Listen(*listenAddr)
	if listenErr != nil {
		return
	}
}

// docker run --name mongodb -d -p 27017:27017 mongo:latest
