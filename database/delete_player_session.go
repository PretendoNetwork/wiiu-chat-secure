package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func DeletePlayerSession(pid uint32) {
	filter := bson.D{
		{"pid", pid},
	}

	_, err := sessionsCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
}
