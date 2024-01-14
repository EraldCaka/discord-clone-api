package db

import (
	"context"
	"errors"
	"github.com/EraldCaka/discord-clone-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageStore interface {
	GetMessageByID(context.Context, string) (*types.Message, error)
	GetMessages(context.Context) ([]*types.Message, error)
	CreateMessage(context.Context, *types.Message) (*types.Message, error)
	DeleteMessage(ctx context.Context, messageID string) error
	//UpdateChannel(ctx context.Context, filter bson.M, params types.UpdateChannelParams) error
}

type MongoMessageStore struct {
	client *mongo.Client
	coll   *mongo.Collection
	ChannelStore
	UserStore
}

func NewMongoMessageStore(client *mongo.Client, channelStore ChannelStore, userStore UserStore) *MongoMessageStore {
	return &MongoMessageStore{
		client:       client,
		coll:         client.Database(DBNAME).Collection(MESSAGE),
		ChannelStore: channelStore,
		UserStore:    userStore,
	}
}

func (s *MongoMessageStore) GetMessageByID(ctx context.Context, id string) (*types.Message, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var message types.Message
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&message); err != nil {
		return nil, err
	}
	return &message, nil
}
func (s *MongoMessageStore) CreateMessage(ctx context.Context, message *types.Message) (*types.Message, error) {
	_, channelDoesntExist := s.ChannelStore.GetChannelByID(ctx, message.ChannelID.Hex())
	if channelDoesntExist != nil {
		return nil, errors.New("channel not found")
	}
	_, userDoesntExist := s.UserStore.GetUserByID(ctx, message.UserID.Hex())
	if userDoesntExist != nil {
		return nil, errors.New("user not found")
	}

	res, err := s.coll.InsertOne(ctx, message)
	if err != nil {
		return nil, err
	}
	message.ID = res.InsertedID.(primitive.ObjectID)
	filter := bson.M{"_id": message.ChannelID}
	update := bson.M{"$push": bson.M{"messages": res.InsertedID}}
	if err := s.ChannelStore.Update(ctx, filter, update); err != nil {
		return nil, err
	}
	return message, nil
}

func (s *MongoMessageStore) GetMessages(ctx context.Context) ([]*types.Message, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var messages []*types.Message
	if err := cur.All(ctx, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}
func (s *MongoMessageStore) DeleteMessage(ctx context.Context, id string) error {
	messageID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := s.GetMessageByID(ctx, id)
	if err := s.ChannelStore.Delete(ctx, res.ChannelID, messageID); err != nil {
		return err
	}
	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": messageID})
	if err != nil {
		return err
	}
	return nil
}
