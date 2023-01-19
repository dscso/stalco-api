package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"rest-go/db"
	"rest-go/models"
	"rest-go/util"
)

type CreateSensorResponse struct {
	Status string             `json:"status"`
	Data   models.SensorModel `json:"data"`
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
		if mongo.IsDuplicateKeyError(err) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "error", "message": "Sensor already exists"})
		}
		return db.ErrorHandler(err)
	}

	return c.JSON(CreateSensorResponse{Status: "success", Data: sensor})
}

type GetSensorResponse struct {
	Status string               `json:"status"`
	Data   []models.SensorModel `json:"data"`
}

// Get sensors
// @Router /area/:area_id/sensors [get]
// @Description get all sensors
// @Response 200 {object} GetSensorResponse
// @Security ApiKeyAuth
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
	return c.JSON(GetSensorResponse{Status: "success", Data: sensors})
}

type GetSensorByAreaResponse struct {
	Status string              `json:"status"`
	Data   []models.SensorData `json:"data"`
}

func GetSensorByArea(c *fiber.Ctx) error {
	area, err := db.GetArea(c)
	if err != nil {
		return err
	}

	filter := bson.D{{"area", area.ID}}
	cursor, err := db.SensorCollection.Find(c.Context(), filter)
	// get array of all sensor ids
	var sensors []models.SensorModel
	if err = cursor.All(c.Context(), &sensors); err != nil {
		return db.ErrorHandler(err)
	}
	var sensorIds []primitive.ObjectID
	for _, sensor := range sensors {
		sensorIds = append(sensorIds, sensor.ID)
	}
	// filter for all sensor data with sensor ids and return the latest data
	filter = bson.D{{"sensor_id", bson.D{{"$in", sensorIds}}}}
	cursor, err = db.SensorDataCollection.Find(c.Context(), filter)
	if err != nil {
		return db.ErrorHandler(err)
	}
	var sensorData []models.SensorData
	if err = cursor.All(c.Context(), &sensorData); err != nil {
		return db.ErrorHandler(err)
	}
	return c.JSON(GetSensorByAreaResponse{Status: "success", Data: sensorData})
}
