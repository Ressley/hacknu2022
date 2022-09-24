package services

import (
	"context"
	"time"

	"github.com/Ressley/hacknu/internal/app/apiserver/helpers"
	"github.com/Ressley/hacknu/internal/app/apiserver/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var apartmentCollection *mongo.Collection = client.Database(helpers.DB).Collection(helpers.APARTMENT)

func CreateApartment(apartment *models.Apartment, building_id *primitive.ObjectID) error {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	apartment.ID = primitive.NewObjectID()
	apartment.Building_id = building_id.Hex()
	_, err = apartmentCollection.InsertOne(ctx, apartment)
	return nil
}

func GetApartment(id *string) (models.Apartment, error) {
	apartment, err := GetApartmentByID(id)
	if err != nil {
		return models.Apartment{}, err
	}
	return apartment, nil
}

func GetApartmentByID(id *string) (models.Apartment, error) {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	result := models.Apartment{}
	ID, _ := primitive.ObjectIDFromHex(*id)
	filter := bson.D{{Key: "_id", Value: ID}}

	err = apartmentCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func AppendApartmentPhoto(apartment *models.Apartment, fileId *string) error {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	var upd bson.D

	filter := bson.D{{Key: "_id", Value: apartment.ID}}

	apartment.Photo = append(apartment.Photo, *fileId)

	upd = bson.D{
		primitive.E{Key: "photo", Value: apartment.Photo},
	}
	updater := bson.D{primitive.E{Key: "$set", Value: upd}}
	_, err = apartmentCollection.UpdateOne(ctx, filter, updater)
	if err != nil {
		return err
	}
	return nil
}
