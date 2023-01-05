package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rest-go/db"
	"rest-go/middleware"
	"rest-go/models"
	"rest-go/util"
)

type CreateSensorResponse struct {
	Status string
	Data   models.SensorModel
}

// CreateUser Creates a user
// @Router /sensors/ [post]
// @Param User body models.SensorModel
// @Description create new sensor
// @Response 200 {object} CreateSensorResponse
func CreateSensor(c *fiber.Ctx) error {
	user := c.Locals("session").(*middleware.SessionAuthenticated)
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
