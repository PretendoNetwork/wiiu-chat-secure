package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func GetAllSessionPIDs() []uint32 {
	var result []bson.M
	output := []uint32{}

	c, _ := sessionsCollection.Find(context.TODO(), bson.D{})
	c.All(context.TODO(), &result)
	for _, i := range result {
		output = append(output, uint32(i["pid"].(int64)))
	}
	return output
}
