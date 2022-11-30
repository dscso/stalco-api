package db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var AreasCollection *mongo.Collection

func init() {
	host := os.Getenv("MONGO_HOST")
	dbName := os.Getenv("MONGO_DB_NAME")

	println("Connecting to host: " + host + " and db: " + dbName + "...")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(host))
	if err != nil {
		panic(err)
	}
	DB = client.Database(dbName)
	AreasCollection = DB.Collection("areas")
}
