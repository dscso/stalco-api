package db

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rest-go/models"
	"rest-go/util"
)

func GetArea(c *fiber.Ctx) (*models.Area, error) {
	idString := c.AllParams()["area_id"]

	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return nil, &fiber.Error{Message: "Invalid ID", Code: fiber.StatusBadRequest}
	}

	var area models.Area

	err = AreasCollection.FindOne(c.Context(), bson.M{"_id": id}).Decode(&area)

	if err != nil {
		return nil, ErrorHandler(err)
	}

	return &area, nil
}

func GetFloor(c *fiber.Ctx, area *models.Area) (*models.Floor, error) {
	floorIdString := c.AllParams()["floor_id"]
	floorId, err := primitive.ObjectIDFromHex(floorIdString)
	if err != nil {
		return nil, &fiber.Error{Message: "Invalid ID", Code: fiber.StatusBadRequest}
	}

	var floor models.Floor
	found := false
	// find floor in area
	for _, v := range area.Floors {
		if v.ID == floorId {
			floor = v
			found = true
		}
	}
	if !found {
		return nil, &fiber.Error{Message: "Floor not found", Code: fiber.StatusNotFound}
	}
	return &floor, nil
}

func GetZone(c *fiber.Ctx, floor *models.Floor) (*models.Zone, error) {
	zoneIdString := c.AllParams()["zone_id"]
	zoneId, err := primitive.ObjectIDFromHex(zoneIdString)
	if err != nil {
		return nil, &fiber.Error{Message: "Invalid ID", Code: fiber.StatusBadRequest}
	}

	var zone models.Zone
	found := false
	// find floor in area
	for _, v := range floor.Zones {
		if v.ID == zoneId {
			zone = v
			found = true
		}
	}
	if !found {
		return nil, &fiber.Error{Message: "Zone not found", Code: fiber.StatusNotFound}
	}
	return &zone, nil
}

// fetches area from database and checks if user is an administrator
func GetAreaForUser(c *fiber.Ctx) (*models.Area, error) {
	session := c.Locals("session").(*util.SessionAuthenticated)
	areaIdString := c.AllParams()["area_id"]

	areaId, err := primitive.ObjectIDFromHex(areaIdString)
	if err != nil {
		return nil, &fiber.Error{Message: "Invalid ID", Code: fiber.StatusBadRequest}
	}

	var area models.Area

	err = AreasCollection.FindOne(c.Context(), bson.M{"_id": areaId}).Decode(&area)
	if err != nil {
		return nil, ErrorHandler(err)
	}

	if !util.ContainsObjectID(area.Administrators, session.UserID) {
		return nil, &fiber.Error{Message: "You are not an administrator of this area", Code: fiber.StatusUnauthorized}
	}
	return &area, nil
}
