package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func EndCall(caller uint32) {
	filter := bson.D{
		{"caller_pid", caller},
	}

	_, err := callsCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
}
