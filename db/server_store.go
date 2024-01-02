package db

import (
	"context"
	"github.com/EraldCaka/discord-clone-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const serverColl = "server"

type ServerStore interface {
	GetServerByID(context.Context, string) (*types.Server, error)
	GetServers(context.Context) ([]*types.Server, error)
	CreateServer(context.Context, *types.Server) (*types.Server, error)
	DeleteServer(context.Context, string) error
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
		coll:      client.Database(DBNAME).Collection(serverColl),
		UserStore: userStore,
	}
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

func (s *MongoServerStore) DeleteServer(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	return nil
}
func (s *MongoServerStore) CreateServer(ctx context.Context, server *types.Server) (*types.Server, error) {
	res, err := s.coll.InsertOne(ctx, server)
	if err != nil {
		return nil, err
	}
	server.ID = res.InsertedID.(primitive.ObjectID)
	filter := bson.M{"_id": server.UserID} // verification if the user exists
	update := bson.M{"$push": bson.M{"server": server.ID}}
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
