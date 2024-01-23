package db

import (
	"context"
	"fmt"
	"github.com/EraldCaka/discord-clone-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Dropper interface {
	Drop(context.Context) error
}

type UserStore interface {
	Dropper
	GetUserByID(ctx context.Context, id string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	CreateUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(ctx context.Context, client *mongo.Client, server []*types.Server, user *types.User) error
	UpdateUser(ctx context.Context, filter bson.M, params types.UpdateUserParams) error
	GetUserByEmail(ctx context.Context, email string) (*types.User, error)
	Update(ctx context.Context, filter bson.M, update bson.M) error
	Delete(ctx context.Context, userID primitive.ObjectID, serverID primitive.ObjectID) error
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(NAME).Collection(USER),
	}
}
func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("--- dropping user collection")
	return s.coll.Drop(ctx)
}
func (s *MongoUserStore) Update(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}

func (s *MongoUserStore) Delete(ctx context.Context, userID primitive.ObjectID, serverID primitive.ObjectID) error {
	filter := bson.M{"_id": userID}
	update := bson.M{"$pull": bson.M{"ownedServers": serverID}}
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}
func (s *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	//fmt.Println(&user)
	return &user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {

		}
	}(cur, ctx)

	var users []*types.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *MongoUserStore) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, client *mongo.Client, server []*types.Server, user *types.User) error {
	for _, serverID := range user.OwnedServers {
		for _, serverObj := range server {
			if serverObj.ID.Hex() == serverID.Hex() {
				if _, err := client.Database(NAME).Collection(MESSAGE).DeleteMany(
					ctx,
					bson.M{"channelID": bson.M{"$in": serverObj.Channels}},
				); err != nil {
					return err
				}
				if _, err := client.Database(NAME).Collection(CHANNEL).DeleteMany(
					ctx,
					bson.M{"serverID": serverObj.ID},
				); err != nil {
					return err
				}
				if _, err := client.Database(NAME).Collection(SERVER).DeleteOne(
					ctx,
					bson.M{"_id": serverObj.ID},
				); err != nil {
					return err
				}
				fmt.Println("Deleted server:", serverObj.ID.Hex())
				break
			}
		}
	}
	userFilter := bson.M{"_id": user.ID}
	if _, err := s.coll.DeleteOne(ctx, userFilter); err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter bson.M, params types.UpdateUserParams) error {
	update := bson.D{
		{
			"$set", params.ToBSON(),
		},
	}
	_, err := s.coll.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
