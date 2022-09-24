package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Apartment struct {
	ID          primitive.ObjectID `bson:"_id, omitempty"`
	Building_id string             `bson:"building_id"`
	Rooms       int                `json:"rooms" validate:"required"`
	Price       int                `json:"price" validate:"required"`
	Count       int                `json:"count" validate:"required"`
	Area        int                `json:"area" validate:"required"`
	Photo       []PhotoData        `json:"photo"`
}

type ApartmentResponse struct {
	Id          string      `json:"id"`
	Building_id string      `json:"building_id"`
	Rooms       int         `json:"rooms"`
	Price       int         `json:"price"`
	Count       int         `json:"count"`
	Area        int         `json:"area"`
	Photo       []PhotoData `json:"photo"`
}
