package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateActivityParams struct {
	Name      string              `json:"name"`
	Details   string              `json:"details"`
	State     string              `json:"state"`
	TimeStamp primitive.Timestamp `json:"timeStamp"`
}

type Activity struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string              `bson:"name" json:"name"`
	Details   string              `bson:"details" json:"details"`
	State     string              `bson:"state" json:"state"`
	TimeStamp primitive.Timestamp `bson:"timeStamp" json:"timeStamp"`
}

func NewActivity(params CreateActivityParams) (*Activity, error) {
	return &Activity{
		Name:      params.Name,
		Details:   params.Details,
		State:     params.State,
		TimeStamp: params.TimeStamp,
	}, nil
}
