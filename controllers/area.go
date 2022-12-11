package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rest-go/db"
	"rest-go/middleware"
	"rest-go/models"
	"rest-go/util"
	"time"
)

type GetAreaResponse struct {
	Status string      `json:"status"`
	Data   models.Area `json:"data"`
}

// GetArea gets an area by id
// @Router /api/area/{area_id} [get]
// @Param area_id path string true "Area ID"
// @Description Get an area by id
// @Response 200 {object} GetAreaResponse
func GetArea(c *fiber.Ctx) error {
	idString := c.AllParams()["area_id"]

	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return &fiber.Error{Message: "Invalid ID", Code: fiber.StatusBadRequest}
	}

	var ar models.Area

	err = db.AreasCollection.FindOne(c.Context(), bson.M{"_id": id}).Decode(&ar)

	if err != nil {
		if err != nil {
			return db.ErrorHandler(err)
		}
	}

	return c.JSON(GetAreaResponse{Status: "success", Data: ar})
}

type EditFloorResponse struct {
	Status string       `json:"status"`
	Data   models.Floor `json:"data"`
}

func EditFloor(c *fiber.Ctx) error {
	session := c.Locals("session").(*middleware.SessionAuthenticated)
	areaIdString := c.AllParams()["area_id"]
	floorIdString := c.AllParams()["floor_id"]

	areaId, err := primitive.ObjectIDFromHex(areaIdString)
	floorId, err2 := primitive.ObjectIDFromHex(floorIdString)

	if err != nil || err2 != nil {
		return &fiber.Error{Message: "Invalid ID", Code: fiber.StatusBadRequest}
	}

	var ar models.Area
	err = db.AreasCollection.FindOne(c.Context(), bson.M{"_id": areaId}).Decode(&ar)
	if err != nil {
		return db.ErrorHandler(err)
	}

	if !util.ContainsObjectID(ar.Administrators, session.UserID) {
		return &fiber.Error{Message: "You are not an administrator of this area", Code: fiber.StatusUnauthorized}
	}

	var floor models.Floor

	for _, v := range ar.Floors {
		if v.ID == floorId {
			floor = v
		}
	}

	if err := c.BodyParser(&floor); err != nil {
		return err
	}

	newFloorBson := db.ConvertStructToBsonM(&floor, "updateble", "floors.$.")

	filter := bson.D{{"_id", areaId}, {"floors._id", floorId}}
	newFloorBson["updated_at"] = time.Now()
	update := bson.D{{"$set", newFloorBson}}
	result, err := db.AreasCollection.UpdateOne(c.Context(), filter, update)

	if err != nil {
		return db.ErrorHandler(err)
	}
	if result.MatchedCount == 0 {
		return &fiber.Error{Code: fiber.StatusNotFound, Message: "floor not found"}
	}

	return c.JSON(EditFloorResponse{Status: "success", Data: floor})
}
