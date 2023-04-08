package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUsernameFromPID(pid uint32) string {
	var result bson.M
	filter := bson.D{
		{
			Key:   "pid",
			Value: pid,
		},
	}

	err := pnidCollection.FindOne(context.TODO(), filter, options.FindOne()).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ""
		}

		panic(err)
	}

	return result["username"].(string)
}
