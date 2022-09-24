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

var eventCollection *mongo.Collection = client.Database(helpers.DB).Collection(helpers.EVENT)

func CreateEvent(event *models.Event) (*mongo.InsertOneResult, error) {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)

	_, err = GetEventByName(event.Name, event.Admin)
	if err == nil {
		return nil, errors.New("community allready exist")
	}

	id, err := eventCollection.InsertOne(ctx, event)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func GetEventAll() ([]models.Event, error) {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	result := []models.Event{}
	filter := bson.D{}
	cursor, err := eventCollection.Find(ctx, filter)

	err = cursor.All(ctx, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func GetEventByName(name *string, admin *string) (models.Event, error) {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	result := models.Event{}
	filter := bson.D{{Key: "name", Value: name}, {Key: "admin", Value: admin}}

	err = eventCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func GetEventByID(id *string) (models.Event, error) {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	result := models.Event{}
	ID, _ := primitive.ObjectIDFromHex(*id)
	filter := bson.D{{Key: "_id", Value: ID}}

	err = eventCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func DeleteEvent(id string) error {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	ID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: ID}}

	_, err := eventCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func GetCommunityByFilter(name *string, account *models.Account) ([]models.Event, error) {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	result := []models.Event{}

	user, err := GetUserOneByUserID(&account.User_id)
	if err != nil {
		return result, err
	}

	if len(user.Followed) == 0 {
		return result, nil
	}

	filter := bson.M{"$and": []bson.M{bson.M{"admin": bson.M{"$in": user.Followed}}, bson.M{"type.name": name}}}
	cursor, err := eventCollection.Find(ctx, filter)
	if err != nil {
		return result, err
	}
	err = cursor.All(ctx, &result)
	if err != nil {
		return result, err
	}
	return result, nil

}
