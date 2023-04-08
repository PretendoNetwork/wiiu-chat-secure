package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetPlayerURLs(pid uint32) []string {
	var result bson.M
	filter := bson.D{
		{"pid", pid},
	}

	err := sessionsCollection.FindOne(context.TODO(), filter, options.FindOne()).Decode(&result)
	if err != nil {
		panic(err)
	}

	oldurlArray := result["urls"].(bson.A)
	newurlArray := make([]string, len(oldurlArray))
	for i := 0; i < len(oldurlArray); i++ {
		newurlArray[i] = oldurlArray[i].(string)
	}

	return newurlArray
}
