package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newCall(caller uint32, target uint32) {
	_, err := callsCollection.InsertOne(context.TODO(), bson.D{{"caller_pid", caller}, {"target_pid", target}, {"ringing", true}})
	if err != nil {
		panic(err)
	}
}

func getCallInfoByCaller(caller uint32) (uint32, uint32, bool) { // caller pid, target pid, ringing
	var result bson.M
	filter := bson.D{{"caller_pid", caller}}

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

func getCallInfoByTarget(target uint32) (uint32, uint32, bool) { // caller pid, target pid, ringing
	var result bson.M
	filter := bson.D{{"target_pid", target}}

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

func endCallRinging(caller uint32) {
	_, err := callsCollection.UpdateOne(context.TODO(), bson.D{{"caller_pid", caller}}, bson.D{{"$set", bson.D{{"ringing", false}}}})
	if err != nil {
		panic(err)
	}
}

func endCall(caller uint32) {
	_, err := callsCollection.DeleteOne(context.TODO(), bson.D{{"caller_pid", caller}})
	if err != nil {
		panic(err)
	}
}
