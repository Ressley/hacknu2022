package services

import (
	"context"
	"fmt"
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

func GetBuilding(id *string) (models.BuildingResponse, error) {
	building, err := GetBuildingByID(id)
	if err != nil {
		return models.BuildingResponse{}, err
	}
	response := models.BuildingResponse{
		ID:           building.ID,
		Name:         building.Name,
		Latitude:     building.Latitude,
		Longitude:    building.Longitude,
		Started_at:   building.Started_at,
		Ends_at:      building.Ends_at,
		Floors:       building.Floors,
		Neighborhood: building.Neighborhood,
		Photo:        building.Photo,
	}

	minPrice := 2147483647
	minArea := 2147483647
	maxArea := 0
	for i := range building.Apartments {
		apartment, err := GetApartment(&building.Apartments[i])
		if err != nil {
			continue
		}

		apartmentResponse := models.ApartmentResponse{
			Id:          apartment.ID.Hex(),
			Building_id: apartment.Building_id,
			Rooms:       apartment.Rooms,
			Price:       apartment.Price,
			Count:       apartment.Count,
			Area:        apartment.Area,
			Photo:       apartment.Photo,
		}

		for len(response.Apartments) < apartmentResponse.Rooms {
			var apartmentTypes []models.ApartmentResponse
			response.Apartments = append(response.Apartments, apartmentTypes)
		}
		response.Apartments[apartmentResponse.Rooms-1] = append(response.Apartments[apartmentResponse.Rooms-1], apartmentResponse)
		// response.Apartments = append(response.Apartments, apartmentResponse)
		minArea = min(minArea, apartment.Area)
		maxArea = max(maxArea, apartment.Area)
		minPrice = min(minPrice, apartment.Price)
	}

	response.MinPrice = &minPrice
	if minArea == maxArea {
		response.Area = fmt.Sprint(maxArea) + " sq. m."
	} else {
		response.Area = fmt.Sprint(minArea) + " sq. m. - " + fmt.Sprint(maxArea) + " sq. m."
	}
	return response, nil
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

func AppendApartment(apartment_id *primitive.ObjectID, building *models.Building) error {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	var upd bson.D

	filter := bson.D{{Key: "_id", Value: building.ID}}

	building.Apartments = append(building.Apartments, apartment_id.Hex())

	upd = bson.D{
		primitive.E{Key: "apartments", Value: building.Apartments},
	}
	updater := bson.D{primitive.E{Key: "$set", Value: upd}}
	_, err = buildingCollection.UpdateOne(ctx, filter, updater)
	if err != nil {
		return err
	}
	return nil
}

func AppendBuildingPhoto(building *models.Building, fileId *string, _type *string) error {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	var upd bson.D

	filter := bson.D{{Key: "_id", Value: building.ID}}
	link := "http://" + helpers.HOST + ":8080/download/photo?fileid=" + fmt.Sprint(*fileId)
	// photo := models.PhotoData{
	// 	Type: _type,
	// 	Link: &link,
	// }
	building.Photo = append(building.Photo, link)

	upd = bson.D{
		primitive.E{Key: "photo", Value: building.Photo},
	}
	updater := bson.D{primitive.E{Key: "$set", Value: upd}}
	_, err = buildingCollection.UpdateOne(ctx, filter, updater)
	if err != nil {
		return err
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
