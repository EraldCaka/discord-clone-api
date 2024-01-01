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

const dburi = "mongodb://localhost:27017"

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New(config)
	apiv1 := app.Group("/discord/api/v1")

	var (
		userHandler = handlers.NewUserHandler(db.NewMongoUserStore(client))
	)
	//user requests

	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	//server requests

	listenErr := app.Listen(*listenAddr)
	if listenErr != nil {
		return
	}
}

// docker run --name mongodb -d -p 27017:27017 mongo:latest
