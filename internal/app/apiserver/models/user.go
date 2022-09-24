package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id, omitempty"`
	User_id    string             `bson:"user_id, omitempty"`
	First_name *string            `json:"first_name" validate:"required,min=2,max=100"`
	Last_name  *string            `json:"last_name" validate:"required,min=2,max=100"`
	Number     *string            `bson:"number, omitempty"`
	Photo      *string            `bson:"photo, omitempty"`
	Events     []string           `bson:"events, omitempty"`
	Followers  []string           `bson:"followers, omitempty"`
	Followed   []string           `bson:"followed, omitempty"`
}

//City       *string            `bson:"city, omitempty"`
//Communities []string           `bson:"community, omitempty"`
