package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var DB *mongo.Database
var AreasCollection *mongo.Collection
var UsersCollection *mongo.Collection

func init() {
	host := os.Getenv("MONGO_HOST")
	dbName := os.Getenv("MONGO_DB_NAME")

	log.Println("Connecting to host: " + host + " and db: " + dbName + "...")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(host))
	if err != nil {
		panic(err)
	}
	DB = client.Database(dbName)

	AreasCollection = DB.Collection("areas")
	UsersCollection = DB.Collection("users")

	model := mongo.IndexModel{
		Keys: bson.M{
			"email": 1,
		},
		Options: options.Index().SetUnique(true),
	}
	_, err = UsersCollection.Indexes().CreateOne(context.TODO(), model)
	if err != nil {
		log.Println(err.Error())
	}
}
