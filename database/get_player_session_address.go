package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetPlayerSessionAddress(pid uint32) string {
	var result bson.M
	filter := bson.D{
		{"pid", pid},
	}

	fmt.Println(pid)

	err := sessionsCollection.FindOne(context.TODO(), filter, options.FindOne()).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "127.0.0.1:9999"
		} else {
			panic(err)
		}
	}

	return result["ip"].(string) + ":" + result["port"].(string)
}
