package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCallInfoByCaller(caller uint32) (uint32, uint32, bool) { // caller pid, target pid, ringing
	var result bson.M
	filter := bson.D{
		{"caller_pid", caller},
	}

	err := callsCollection.FindOne(context.TODO(), filter, options.FindOne()).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, 0, false
		} else {
			panic(err)
		}
	} else {
		return uint32(result["caller_pid"].(int64)), uint32(result["target_pid"].(int64)), result["ringing"].(bool)
	}
}
