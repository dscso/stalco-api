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

// fetches area from database and checks if user is an administrator
func getAreaForUser(c *fiber.Ctx) (*models.Area, error) {
	session := c.Locals("session").(*middleware.SessionAuthenticated)
	areaIdString := c.AllParams()["area_id"]

	areaId, err := primitive.ObjectIDFromHex(areaIdString)
	if err != nil {
		return nil, &fiber.Error{Message: "Invalid ID", Code: fiber.StatusBadRequest}
	}

	var area models.Area

	err = db.AreasCollection.FindOne(c.Context(), bson.M{"_id": areaId}).Decode(&area)
	if err != nil {
		return nil, db.ErrorHandler(err)
	}

	if !util.ContainsObjectID(area.Administrators, session.UserID) {
		return nil, &fiber.Error{Message: "You are not an administrator of this area", Code: fiber.StatusUnauthorized}
	}
	return &area, nil
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

	var area models.Area

	err = db.AreasCollection.FindOne(c.Context(), bson.M{"_id": id}).Decode(&area)

	if err != nil {
		return db.ErrorHandler(err)
	}

	return c.JSON(GetAreaResponse{Status: "success", Data: area})
}

type EditAreaResponse struct {
	Status string      `json:"status"`
	Data   models.Area `json:"data"`
}

// EditArea edits an area by id
// @Router /api/area/{area_id} [put]
// @Param area_id path string true "Area ID"
// @Description Edit an area by id
// @Response 200 {object} EditAreaResponse
// @Security ApiKeyAuth
func EditArea(c *fiber.Ctx) error {
	areaInDB, err := getAreaForUser(c)
	if err != nil {
		return err
	}

	area := *areaInDB // for modifying the area
	if err := c.BodyParser(&area); err != nil {
		return err
	}

	newAreaJson := db.ConvertStructToBsonM(&area, "updateble", "")

	filter := bson.D{{"_id", areaInDB.ID}}
	newAreaJson["updated_at"] = time.Now()
	update := bson.D{{"$set", newAreaJson}}

	result, err := db.AreasCollection.UpdateOne(c.Context(), filter, update)

	if err != nil {
		return db.ErrorHandler(err)
	}
	if result.MatchedCount == 0 {
		return &fiber.Error{Code: fiber.StatusNotFound, Message: "area not found"}
	}
	// read again from database to get the updated data
	return GetArea(c)
}

type EditFloorResponse struct {
	Status string       `json:"status"`
	Data   models.Floor `json:"data"`
}

func EditFloor(c *fiber.Ctx) error {
	areaInDB, err := getAreaForUser(c)
	if err != nil {
		return err
	}

	floorIdString := c.AllParams()["floor_id"]
	floorId, err := primitive.ObjectIDFromHex(floorIdString)
	if err != nil {
		return &fiber.Error{Message: "Invalid ID", Code: fiber.StatusBadRequest}
	}

	var floor models.Floor
	// find floor in area
	for _, v := range areaInDB.Floors {
		if v.ID == floorId {
			floor = v
		}
	}

	if err := c.BodyParser(&floor); err != nil {
		return err
	}

	newFloorBson := db.ConvertStructToBsonM(&floor, "updateble", "floors.$.")

	filter := bson.D{{"_id", areaInDB.ID}, {"floors._id", floorId}}
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
