package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func AddNewUser(pid uint32) {
	username := GetUsernameFromPID(pid)

	filter := bson.D{
		{"pid", pid},
		{"missed_calls", bson.A{""}},
		{"username", username},
		{"status", "unallowed"},
	}

	_, err := usersCollection.InsertOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
}
