package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Apartment struct {
	ID          primitive.ObjectID `bson:"_id, omitempty"`
	Building_id string             `bson:"building_id"`
	Rooms       int                `json:"rooms" validate:"required"`
	Price       int                `json:"price" validate:"required"`
	Count       int                `json:"count" validate:"required"`
	Area        int                `json:"area" validate:"required"`
	Photo       []string           `json:"photo"`
}
