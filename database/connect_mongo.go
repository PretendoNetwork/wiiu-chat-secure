package database

import (
	"context"
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
	mongoClient, _ = mongo.NewClient(options.Client().ApplyURI(os.Getenv("PN_WIIU_CHAT_MONGO_URI")))
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
