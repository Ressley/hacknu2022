package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Event struct {
	ID           primitive.ObjectID `bson:"_id, omitempty"`
	event_id     *string            `bson:"event_id, "`
	Name         *string            `bson:"name, omitempty"`
	Photo        *string            `bson:"photo, omitempty"`
	City         *string            `bson:"city: omitempty"`
	Type         *Community         `bson:"type, omitempty"`
	TypeName     *string            `bson:"typeName, omitempty"`
	Description  *string            `bson:"description, omitempty"`
	Participants []string           `bson:"participants, omitempty"`
	Admin        *string            `bson:"admin, omitempty"`
	Time         *string            `bson:"time, omitempty"`
}
