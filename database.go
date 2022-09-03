package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var mongoContext context.Context
var accountDatabase *mongo.Database
var doorsDatabase *mongo.Database
var pnidCollection *mongo.Collection
var nexAccountsCollection *mongo.Collection
var regionsCollection *mongo.Collection
var usersCollection *mongo.Collection
var sessionsCollection *mongo.Collection
var callsCollection *mongo.Collection
var tourneysCollection *mongo.Collection

func connectMongo() {
	mongoClient, _ = mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	mongoContext, _ = context.WithTimeout(context.Background(), 10*time.Second)
	_ = mongoClient.Connect(mongoContext)

	accountDatabase = mongoClient.Database("pretendo")
	pnidCollection = accountDatabase.Collection("pnids")
	nexAccountsCollection = accountDatabase.Collection("nexaccounts")

	doorsDatabase = mongoClient.Database("doors")
	usersCollection = doorsDatabase.Collection("users")
	sessionsCollection = doorsDatabase.Collection("sessions")
	callsCollection = doorsDatabase.Collection("calls")

	sessionsCollection.DeleteMany(context.TODO(), bson.D{})
	callsCollection.DeleteMany(context.TODO(), bson.D{})
}

func getUsernameFromPID(pid uint32) string {
	var result bson.M

	err := pnidCollection.FindOne(context.TODO(), bson.D{{Key: "pid", Value: pid}}, options.FindOne()).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ""
		}

		panic(err)
	}

	return result["username"].(string)
}

func addNewUser(pid uint32) {
	_, err := usersCollection.InsertOne(context.TODO(), bson.D{{"pid", pid}, {"missed_calls", bson.A{""}}, {"username", getUsernameFromPID(pid)}, {"status", "unallowed"}})
	if err != nil {
		panic(err)
	}
}

func isUserAllowed(pid uint32) bool {
	var result bson.M

	err := usersCollection.FindOne(context.TODO(), bson.D{{"pid", pid}}, options.FindOne()).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		} else {
			panic(err)
		}
	} else {
		data := result["status"].(string)
		if data == "allowed" {
			return true
		} else {
			return false
		}
	}
}

func doesUserExist(pid uint32) bool {
	var result bson.M

	err := usersCollection.FindOne(context.TODO(), bson.D{{"pid", pid}}, options.FindOne()).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		} else {
			panic(err)
		}
	} else {
		return true
	}
}

func addPlayerSession(pid uint32, urls []string, ip string, port string) {
	_, err := sessionsCollection.InsertOne(context.TODO(), bson.D{{"pid", pid}, {"urls", urls}, {"ip", ip}, {"port", port}})
	if err != nil {
		panic(err)
	}
}

func getAllSessionPIDs() []uint32 {
	var result []bson.M
	output := []uint32{}

	c, _ := sessionsCollection.Find(context.TODO(), bson.D{})
	c.All(context.TODO(), &result)
	for _, i := range result {
		output = append(output, uint32(i["pid"].(int64)))
	}
	return output
}

func doesSessionExist(pid uint32) bool {
	var result bson.M

	err := sessionsCollection.FindOne(context.TODO(), bson.D{{"pid", pid}}, options.FindOne()).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		} else {
			panic(err)
		}
	} else {
		return true
	}
}

func updatePlayerSessionAll(pid uint32, urls []string, ip string, port string) {
	_, err := sessionsCollection.UpdateOne(context.TODO(), bson.D{{"pid", pid}}, bson.D{{"$set", bson.D{{"pid", pid}, {"urls", urls}, {"ip", ip}, {"port", port}}}})
	if err != nil {
		panic(err)
	}
}

func updatePlayerSessionUrl(pid uint32, oldurl string, newurl string) {
	var result bson.M

	err := sessionsCollection.FindOne(context.TODO(), bson.D{{"pid", pid}}, options.FindOne()).Decode(&result)
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

	_, err = sessionsCollection.UpdateOne(context.TODO(), bson.D{{"pid", pid}}, bson.D{{"$set", bson.D{{"urls", newurlArray}}}})
	if err != nil {
		panic(err)
	}
}

func deletePlayerSession(pid uint32) {
	_, err := sessionsCollection.DeleteOne(context.TODO(), bson.D{{"pid", pid}})
	if err != nil {
		panic(err)
	}
}

func getPlayerUrls(pid uint32) []string {
	var result bson.M

	err := sessionsCollection.FindOne(context.TODO(), bson.D{{"pid", pid}}, options.FindOne()).Decode(&result)
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

func getPlayerSessionAddress(pid uint32) string {
	var result bson.M

	fmt.Println(pid)

	err := sessionsCollection.FindOne(context.TODO(), bson.D{{"pid", pid}}, options.FindOne()).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "127.0.0.1:9999"
		} else {
			panic(err)
		}
	}

	return result["ip"].(string) + ":" + result["port"].(string)
}
