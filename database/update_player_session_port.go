package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdatePlayerSessionPort(pid uint32, port string) {
	filter := bson.D{
		{"pid", pid},
	}

	update := bson.D{
		{
			"$set", bson.D{
				{"port", port},
			},
		},
	}

	_, err := sessionsCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
}
