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
