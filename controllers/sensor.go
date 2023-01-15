package controllers

import (
	"github.com/gofiber/fiber/v2"
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
	user := c.Locals("session").(*util.SessionAuthenticated)
	if user.Authenticated == false {
		return util.UnauthorizedError
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

	// inserting user in database
	_, err = db.SensorCollection.InsertOne(c.Context(), sensor)
	if err != nil {
		return db.ErrorHandler(err)
	}

	return c.JSON(CreateSensorResponse{Status: "success", Data: sensor})
}
