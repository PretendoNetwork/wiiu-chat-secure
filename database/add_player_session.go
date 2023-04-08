package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func AddPlayerSession(pid uint32, urls []string, ip string, port string) {
	filter := bson.D{
		{"pid", pid},
		{"urls", urls},
		{"ip", ip},
		{"port", port},
	}

	_, err := sessionsCollection.InsertOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
}
