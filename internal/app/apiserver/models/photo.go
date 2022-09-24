package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Photo struct {
	ID     primitive.ObjectID `bson:"_id, omitempty"`
	Fileid *string            `bson:"fileid"`
	Hash   *string            `bson:"hash"`
}
