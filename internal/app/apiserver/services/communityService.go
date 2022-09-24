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

var communityCollection *mongo.Collection = client.Database(helpers.DB).Collection(helpers.COMMUNITY)

func CreateCommunity(community *models.Community) error {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	_, err = GetCommunityByName(community.Name)
	if err == nil {
		return errors.New("community allready exist")
	}
	community.ID = primitive.NewObjectID()
	_, err = communityCollection.InsertOne(ctx, community)
	if err != nil {
		return err
	}
	return nil
}

func GetCommunityByNameAndAdmin(name *string, admin *string) (models.Community, error) {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	result := models.Community{}
	filter := bson.D{{Key: "name", Value: name}, {Key: "admin", Value: admin}}

	err = communityCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil

}

func GetCommunityByName(name *string) (models.Community, error) {
	var ctx, _ = context.WithTimeout(context.TODO(), 100*time.Second)
	result := models.Community{}
	filter := bson.D{{Key: "name", Value: name}}

	err = communityCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil

}
