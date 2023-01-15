package db

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var AreasCollection *mongo.Collection
var UsersCollection *mongo.Collection
var SessionsCollection *mongo.Collection
var SensorCollection *mongo.Collection
var SensorDataCollection *mongo.Collection

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
	SessionsCollection = DB.Collection("sessions")
	SensorCollection = DB.Collection("sensors")
	SensorDataCollection = DB.Collection("sensor_data")

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

func ErrorHandler(err error) error {
	if err == mongo.ErrNoDocuments {
		return &fiber.Error{Message: "Not found in the Database", Code: fiber.StatusNotFound}
	} else {
		log.Println(err.Error())
		return &fiber.Error{Message: "Internal server error", Code: fiber.StatusInternalServerError}
	}
}
