package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DoesUserExist(pid uint32) bool {
	var result bson.M
	filter := bson.D{
		{"pid", pid},
	}

	err := usersCollection.FindOne(context.TODO(), filter, options.FindOne()).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		} else {
			panic(err)
		}
	} else {
		return true
	}
}
