package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdatePlayerSessionUrl(pid uint32, oldurl string, newurl string) {
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
		if oldurlArray[i].(string) == oldurl {
			newurlArray[i] = newurl
		} else {
			newurlArray[i] = oldurlArray[i].(string)
		}
	}

	update := bson.D{
		{
			"$set", bson.D{
				{"urls", newurlArray},
			},
		},
	}

	_, err = sessionsCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
}
