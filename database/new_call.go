package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func NewCall(caller uint32, target uint32) {
	document := bson.D{
		{"caller_pid", caller},
		{"target_pid", target},
		{"ringing", true},
	}

	_, err := callsCollection.InsertOne(context.TODO(), document)
	if err != nil {
		panic(err)
	}
}
