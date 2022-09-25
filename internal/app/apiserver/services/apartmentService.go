package services

import (
	"context"
	"fmt"
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

func AppendApartmentPhoto(apartment *models.Apartment, fileId *string, _type *string, name *string) error {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	var upd bson.D

	filter := bson.D{{Key: "_id", Value: apartment.ID}}
	link := "http://" + helpers.HOST + ":8080/download/photo?fileid=" + fmt.Sprint(*fileId)
	photo := models.PhotoData{
		Type: _type,
		Name: name,
		Link: &link,
	}
	apartment.Photo = append(apartment.Photo, photo)

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

func DeleteApartment(id *primitive.ObjectID) error {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	filter := bson.D{{Key: "_id", Value: id}}

	_, err := apartmentCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func RemoveApartmentPhoto(apartment *models.Apartment, fileId *string) error {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	var upd bson.D

	filter := bson.D{{Key: "_id", Value: apartment.ID}}
	link := "http://" + helpers.HOST + ":8080/download/photo?fileid=" + fmt.Sprint(*fileId)
	for i := range apartment.Photo {
		if *apartment.Photo[i].Link == link {
			apartment.Photo = RemoveIndex(apartment.Photo, i)
			break
		}
	}

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
