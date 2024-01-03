package db

import (
	"context"
	"fmt"
	"github.com/EraldCaka/discord-clone-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChannelStore interface {
	GetChannelByID(context.Context, string) (*types.Channel, error)
	GetChannels(context.Context) ([]*types.Channel, error)
	CreateChannel(context.Context, *types.Channel) (*types.Channel, error)
	DeleteChannel(context.Context, string) error
	//UpdateChannel(ctx context.Context, filter bson.M, params types.UpdateChannelParams) error
}

type MongoChannelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
	ServerStore
}

func NewMongoChannelStore(client *mongo.Client, serverStore ServerStore) *MongoChannelStore {
	return &MongoChannelStore{
		client:      client,
		coll:        client.Database(DBNAME).Collection(Channel),
		ServerStore: serverStore,
	}
}
func (s *MongoChannelStore) GetChannelByID(ctx context.Context, id string) (*types.Channel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var channel types.Channel
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&channel); err != nil {
		return nil, err
	}
	return &channel, nil
}
func (s *MongoChannelStore) CreateChannel(ctx context.Context, channel *types.Channel) (*types.Channel, error) {
	_, serverDoesntExist := s.ServerStore.GetServerByID(ctx, channel.ServerID.Hex())
	fmt.Println(channel.ServerID)
	if serverDoesntExist != nil {
		return nil, serverDoesntExist
	}
	res, err := s.coll.InsertOne(ctx, channel)
	if err != nil {
		return nil, err
	}
	channel.ID = res.InsertedID.(primitive.ObjectID)
	filter := bson.M{"_id": channel.ServerID}
	update := bson.M{"$push": bson.M{"channels": res.InsertedID}}
	if err := s.ServerStore.Update(ctx, filter, update); err != nil {
		return nil, err
	}
	return channel, nil
}

func (s *MongoChannelStore) GetChannels(ctx context.Context) ([]*types.Channel, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var channels []*types.Channel
	if err := cur.All(ctx, &channels); err != nil {
		return nil, err
	}
	return channels, nil
}
func (s *MongoChannelStore) DeleteChannel(ctx context.Context, id string) error {
	channelID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := s.GetChannelByID(ctx, id)
	if err := s.ServerStore.Delete(ctx, res.ServerID, channelID); err != nil {
		return err
	}
	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": channelID})
	if err != nil {
		return err
	}
	return nil
}
