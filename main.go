package main

import (
	"context"
	_ "context"
	"flag"
	"github.com/EraldCaka/discord-clone-api/db"
	"github.com/EraldCaka/discord-clone-api/handlers"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"log"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New(config)
	api := app.Group("/discord/api/v1")

	var (
		userStore   = db.NewMongoUserStore(client)
		serverStore = db.NewMongoServerStore(client, userStore)
		store       = &db.Store{
			User:   userStore,
			Server: serverStore,
		}
		userHandler   = handlers.NewUserHandler(store)
		serverHandler = handlers.NewServerHandler(serverStore)
	)

	api.Get("/user/:id", userHandler.HandleGetUser)
	api.Put("/user/:id", userHandler.HandlePutUser)
	api.Delete("/user/:id", userHandler.HandleDeleteUser)
	api.Post("/user", userHandler.HandlePostUser)
	api.Get("/user", userHandler.HandleGetUsers)

	api.Delete("/server/:id", serverHandler.HandleDeleteServer)
	api.Post("/server", serverHandler.HandleCreateServer)
	api.Get("/server", serverHandler.HandleGetServers)
	api.Get("/server/:id", serverHandler.HandleGetServer)

	listenErr := app.Listen(*listenAddr)
	if listenErr != nil {
		return
	}
}

// docker run --name mongodb -d -p 27017:27017 mongo:latest
