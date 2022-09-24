package helpers

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientInstance *mongo.Client
var clientInstanceError error
var mongoOnce sync.Once

const (
	CONNECTIONSTRING = "mongodb://localhost:27017"
	HOST             = "localhost"
	DB               = "HackNUDB"
	ACCOUNTS         = "Accounts"
	USERS            = "Users"
	COMMUNITY        = "Community"
	BUILDINGS        = "Buildings"
	APARTMENT        = "Apartments"
	EVENT            = "Event"

	PHOTO_DB     = "HackNUPhoto"
	FILE_DB      = "HackNUFile"
	FS_FILES     = "fs.files"
	FS_CHUNKS    = "fs.chunks"
	HASHED_PHOTO = "hashed.photo"
)

func GetMongoClient() (*mongo.Client, error) {
	mongoOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(CONNECTIONSTRING)
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
		}

		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
		}
		clientInstance = client
	})
	return clientInstance, clientInstanceError
}
