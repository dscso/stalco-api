package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"rest-go/db"
	"rest-go/models"
)

type GetAreaResponse struct {
	Status string      `json:"status"`
	Data   models.Area `json:"data"`
}

func GetArea(c *fiber.Ctx) error {
	idString := c.AllParams()["id"]

	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return &fiber.Error{Message: "Invalid ID", Code: fiber.StatusBadRequest}
	}

	var ar models.Area

	err = db.AreasCollection.FindOne(c.Context(), bson.M{"_id": id}).Decode(&ar)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &fiber.Error{Message: "Area not found", Code: fiber.StatusNotFound}
		} else {
			return &fiber.Error{Message: "Internal server error", Code: fiber.StatusInternalServerError}
		}
	}

	return c.JSON(GetAreaResponse{Status: "success", Data: ar})
}
