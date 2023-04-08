package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func EndCallRinging(caller uint32) {
	filter := bson.D{
		{"caller_pid", caller},
	}

	update := bson.D{
		{
			"$set", bson.D{
				{"ringing", false},
			},
		},
	}

	_, err := callsCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
}
