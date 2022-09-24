package services

import (
	"context"
	"time"

	"github.com/Ressley/hacknu/internal/app/apiserver/helpers"
	"github.com/Ressley/hacknu/internal/app/apiserver/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var apartmentCollection *mongo.Collection = client.Database(helpers.DB).Collection(helpers.BUILDINGS)

func CreateApartment(apartment *models.Apartment, building_id *primitive.ObjectID) error {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	apartment.ID = primitive.NewObjectID()
	apartment.Building_id = building_id.Hex()
	_, err = apartmentCollection.InsertOne(ctx, apartment)
	return nil
}
