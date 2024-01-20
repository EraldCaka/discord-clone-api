package db

import (
	"context"
	"github.com/EraldCaka/discord-clone-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type ChannelStore interface {
	GetChannelByID(context.Context, string) (*types.Channel, error)
	GetChannels(context.Context) ([]*types.Channel, error)
	CreateChannel(context.Context, *types.Channel) (*types.Channel, error)
	DeleteChannel(ctx context.Context, id string, client *mongo.Client) error
	Update(ctx context.Context, filter bson.M, update bson.M) error
	Delete(ctx context.Context, channelID primitive.ObjectID, messageID primitive.ObjectID) error

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
		coll:        client.Database(NAME).Collection(CHANNEL),
		ServerStore: serverStore,
	}
}

func (s *MongoChannelStore) Update(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, update) // CHECKING THE CHANNEL VALUES
	return err

}
func (s *MongoChannelStore) Delete(ctx context.Context, channelID primitive.ObjectID, messageID primitive.ObjectID) error {
	filter := bson.M{"_id": channelID}
	update := bson.M{"$pull": bson.M{"messages": messageID}}
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
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
func (s *MongoChannelStore) DeleteChannel(ctx context.Context, id string, client *mongo.Client) error {
	channelID, err := primitive.ObjectIDFromHex(id)
	channel, err := s.GetChannelByID(ctx, id)
	for _, _ = range channel.Messages {
		if _, err := client.Database(NAME).Collection(MESSAGE).DeleteMany(
			ctx,
			bson.M{"channelID": channelID},
		); err != nil {
			log.Fatal(err)
		}
	}
	if err != nil {
		return err
	}
	if err := s.ServerStore.Delete(ctx, channel.ServerID, channelID); err != nil {
		return err
	}
	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": channelID})
	if err != nil {
		return err
	}
	return nil
}
