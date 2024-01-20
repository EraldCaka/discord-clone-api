package db

import (
	"context"
	"fmt"
	"github.com/EraldCaka/discord-clone-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServerStore interface {
	GetServerByID(context.Context, string) (*types.Server, error)
	GetServers(context.Context) ([]*types.Server, error)
	CreateServer(context.Context, *types.Server) (*types.Server, error)
	DeleteServer(ctx context.Context, client *mongo.Client, id string) error
	Update(ctx context.Context, filter bson.M, update bson.M) error
	Delete(ctx context.Context, serverID primitive.ObjectID, channelID primitive.ObjectID) error
	//UpdateServer(ctx context.Context, filter bson.M, params types.UpdateServerParams) error
}

type MongoServerStore struct {
	client *mongo.Client
	coll   *mongo.Collection
	UserStore
}

func NewMongoServerStore(client *mongo.Client, userStore UserStore) *MongoServerStore {
	return &MongoServerStore{
		client:    client,
		coll:      client.Database(NAME).Collection(SERVER),
		UserStore: userStore,
	}
}

func (s *MongoServerStore) Update(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}

func (s *MongoServerStore) Delete(ctx context.Context, serverID primitive.ObjectID, channelID primitive.ObjectID) error {
	filter := bson.M{"_id": serverID}
	update := bson.M{"$pull": bson.M{"channels": channelID}}
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}

func (s *MongoServerStore) GetServerByID(ctx context.Context, id string) (*types.Server, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var server types.Server
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&server); err != nil {
		return nil, err
	}
	return &server, nil
}

//func (s *MongoServerStore) UpdateServer(ctx context.Context, filter bson.M, params types.UpdateServerParams) error {
//	update := bson.D{
//		{
//			"$set", params.ToBSON(),
//		},
//	}
//	_, err := s.coll.UpdateOne(ctx, filter, update)
//	if err != nil {
//		return err
//	}
//	return nil
//}

func (s *MongoServerStore) DeleteServer(ctx context.Context, client *mongo.Client, id string) error {
	serverID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	server, err := s.GetServerByID(ctx, id)
	if err != nil {
		return err
	}
	if _, err := client.Database(NAME).Collection(MESSAGE).DeleteMany(
		ctx,
		bson.M{"channelID": bson.M{"$in": server.Channels}},
	); err != nil {
		return err
	}
	if _, err := client.Database(NAME).Collection(CHANNEL).DeleteMany(
		ctx,
		bson.M{"serverID": server.ID},
	); err != nil {
		return err
	}
	fmt.Println("Deleted server:", server.ID.Hex())
	res, err := s.GetServerByID(ctx, id)
	if err := s.UserStore.Delete(ctx, res.UserID, serverID); err != nil {
		return err
	}
	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": serverID})
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoServerStore) CreateServer(ctx context.Context, server *types.Server) (*types.Server, error) {
	_, userDoesntExist := s.UserStore.GetUserByID(ctx, server.UserID.Hex())
	if userDoesntExist != nil {
		return nil, userDoesntExist
	}
	res, err := s.coll.InsertOne(ctx, server)
	if err != nil {
		return nil, err
	}
	server.ID = res.InsertedID.(primitive.ObjectID)
	filter := bson.M{"_id": server.UserID}
	update := bson.M{"$push": bson.M{"ownedServers": res.InsertedID}}
	//TODO FIX update since it is being directed to user not to channel store
	if err := s.UserStore.Update(ctx, filter, update); err != nil {
		return nil, err
	}
	return server, nil
}

func (s *MongoServerStore) GetServers(ctx context.Context) ([]*types.Server, error) {

	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var servers []*types.Server
	if err := cur.All(ctx, &servers); err != nil {
		return nil, err
	}
	return servers, nil
}
