package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rest-go/db"
	"rest-go/models"
	"rest-go/util"
)

type CreateSensorResponse struct {
	Status string
	Data   models.SensorModel
}

// CreateSensor Creates a sensor
// @Router /area/:area_id/sensors [post]
// @Param Sensor body models.SensorModel true "Sensor"
// @Description create new sensor
// @Response 200 {object} CreateSensorResponse
// @Security ApiKeyAuth
func CreateSensor(c *fiber.Ctx) error {
	area, err := db.GetArea(c)
	if err != nil {
		return err
	}
	// converting json to struct
	var sensor models.SensorModel
	if err := c.BodyParser(&sensor); err != nil {
		return err
	}
	// hashing password
	key, err := util.Create64ByteKey()
	if err != nil {
		return err
	}

	sensor.ID = primitive.NewObjectID()
	sensor.Key = key
	sensor.Area = area.ID

	// inserting user in database
	_, err = db.SensorCollection.InsertOne(c.Context(), sensor)
	if err != nil {
		return db.ErrorHandler(err)
	}

	return c.JSON(CreateSensorResponse{Status: "success", Data: sensor})
}

func GetSensors(c *fiber.Ctx) error {
	area, err := db.GetArea(c)
	if err != nil {
		return err
	}
	// query form sensor collection all sensors who are associated with the area
	filter := bson.D{{"area", area.ID}}
	cursor, err := db.SensorCollection.Find(c.Context(), filter)
	if err != nil {
		return db.ErrorHandler(err)
	}
	var sensors []models.SensorModel
	if err = cursor.All(c.Context(), &sensors); err != nil {
		return db.ErrorHandler(err)
	}
	return c.JSON(sensors)
}
