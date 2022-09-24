package services

import (
	"context"
	"time"

	"github.com/Ressley/hacknu/internal/app/apiserver/helpers"
	"github.com/Ressley/hacknu/internal/app/apiserver/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var buildingCollection *mongo.Collection = client.Database(helpers.DB).Collection(helpers.BUILDINGS)

func CreateBuilding(building *models.Building) error {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	_, err = GetBuildingByName(building.Name)
	if err == nil {
		return errors.New("building allready exist")
	}
	building.ID = primitive.NewObjectID()
	_, err = buildingCollection.InsertOne(ctx, building)
	if err != nil {
		return err
	}
	return nil
}

func GetBuildingByName(name *string) (models.Building, error) {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	result := models.Building{}
	filter := bson.D{{Key: "name", Value: name}}

	err = buildingCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func GetBuildingByID(id *string) (models.Building, error) {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	result := models.Building{}
	ID, _ := primitive.ObjectIDFromHex(*id)
	filter := bson.D{{Key: "_id", Value: ID}}

	err = buildingCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func ListBuildings() ([]models.BuildingMeta, error) {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	result := []models.BuildingMeta{}
	dbData := []models.Building{}
	filter := bson.D{}

	cursor, err := buildingCollection.Find(ctx, filter)
	if err != nil {
		return result, err
	}
	err = cursor.All(ctx, &dbData)

	for i := range dbData {
		meta := models.BuildingMeta{
			ID:        dbData[i].ID.Hex(),
			Latitude:  *dbData[i].Latitude,
			Longitude: *dbData[i].Longitude,
		}
		result = append(result, meta)
	}

	return result, nil
}
