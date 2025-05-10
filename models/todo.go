package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Completed bool               `json:"completed" bson:"completed, omitempty"`
	Body      string             `json:"body" bson:"body, omitempty"`
}
