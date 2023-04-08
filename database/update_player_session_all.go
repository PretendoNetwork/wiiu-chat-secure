package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdatePlayerSessionAll(pid uint32, urls []string, ip string, port string) {
	filter := bson.D{
		{"pid", pid},
	}

	update := bson.D{
		{
			"$set", bson.D{
				{"pid", pid},
				{"urls", urls},
				{"ip", ip},
				{"port", port},
			},
		},
	}

	_, err := sessionsCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
}
