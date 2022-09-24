package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Building struct {
	ID           primitive.ObjectID `bson:"_id, omitempty"`
	Name         *string            `json:"name" validate:"required,min=2,max=100"`
	Latitude     *string            `json:"latitude" validate:"required"`
	Longitude    *string            `json:"longitude" validate:"required"`
	Started_at   time.Time          `json:"started_at" validate:"required"`
	Ends_at      time.Time          `json:"ends_at" validate:"required"`
	Floors       *int               `json:"floors" validate:"required"`
	Neighborhood []string           `json:"neighborhood"`
	Apartments   []string           `json:"apartments"`
	Photo        []PhotoData        `json:"photo"`
}

type BuildingResponse struct {
	ID           primitive.ObjectID    `json:"id"`
	Name         *string               `json:"name"`
	Latitude     *string               `json:"latitude"`
	Longitude    *string               `json:"longitude"`
	Started_at   time.Time             `json:"started_at"`
	Ends_at      time.Time             `json:"ends_at"`
	Floors       *int                  `json:"floors"`
	Area         string                `json:"area"`
	MinPrice     *int                  `json:"min_price"`
	Neighborhood []string              `json:"neighborhood"`
	Photo        []PhotoData           `json:"photo"`
	Apartments   [][]ApartmentResponse `json:"apartments"`
}

type BuildingMeta struct {
	ID        string `json:"id"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}
