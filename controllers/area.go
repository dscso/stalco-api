package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rest-go/db"
	"rest-go/models"
	"time"
)

type ListAreasResponse struct {
	Status string        `json:"status"`
	Data   []models.Area `json:"data"`
}

// Get all areas for frontend listing
// @Router /api/area/list [get]
// @Description Get all areas for frontend listing
// @Response 200 {object} ListAreasResponse
func GetAreas(c *fiber.Ctx) error {
	// only get these fields form mongodb
	projection := bson.D{
		{"name", 1},
		{"description", 1},
		{"image", 1},
		{"created_at", 1},
		{"updated_at", 1},
	}
	opts := options.Find().SetProjection(projection)
	cursor, err := db.AreasCollection.Find(c.Context(), bson.M{}, opts)
	if err != nil {
		return err
	}
	var areas []models.Area
	for cursor.Next(c.Context()) {
		var area models.Area
		err := cursor.Decode(&area)
		if err != nil {
			return err
		}
		areas = append(areas, area)
	}
	return c.JSON(ListAreasResponse{Status: "success", Data: areas})
}

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
	area, err := db.GetArea(c)
	if err != nil {
		return err
	}
	return c.JSON(GetAreaResponse{Status: "success", Data: *area})
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
	areaInDB, err := db.GetAreaForUser(c)
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

type GetFloorResponse struct {
	Status string       `json:"status"`
	Data   models.Floor `json:"data"`
}

func GetFloor(c *fiber.Ctx) error {
	areaInDB, err := db.GetArea(c)
	if err != nil {
		return err
	}
	floor, err := db.GetFloor(c, areaInDB)
	if err != nil {
		return err
	}
	return c.JSON(GetFloorResponse{Status: "success", Data: *floor})
}

type EditFloorResponse struct {
	Status string       `json:"status"`
	Data   models.Floor `json:"data"`
}

// EditFloor edits a floor from an area
// @Router /api/area/{area_id}/floor/{floor_id} [put]
// @Param area_id path string true "Area ID"
// @Param floor_id path string true "Floor ID"
// @Description Edit a floor from an area
// @Response 200 {object} EditFloorResponse
// @Security ApiKeyAuth
func EditFloor(c *fiber.Ctx) error {
	areaInDB, err := db.GetAreaForUser(c)
	if err != nil {
		return err
	}
	floor, err := db.GetFloor(c, areaInDB)
	if err != nil {
		return err
	}
	if err := c.BodyParser(floor); err != nil {
		return err
	}

	newFloorBson := db.ConvertStructToBsonM(floor, "updateble", "floors.$.")

	filter := bson.D{{"_id", areaInDB.ID}, {"floors._id", floor.ID}}
	newFloorBson["updated_at"] = time.Now()
	update := bson.D{{"$set", newFloorBson}}
	result, err := db.AreasCollection.UpdateOne(c.Context(), filter, update)

	if err != nil {
		return db.ErrorHandler(err)
	}
	if result.MatchedCount == 0 {
		return &fiber.Error{Code: fiber.StatusNotFound, Message: "floor not found"}
	}

	return c.JSON(EditFloorResponse{Status: "success", Data: *floor})
}

type GetZoneResponse struct {
	Status string      `json:"status"`
	Data   models.Zone `json:"data"`
}

// GetZone gets a zone from an area
// @Router /api/area/{area_id}/floor/{floor_id}/zone/{zone_id} [get]
// @Param area_id path string true "Area ID"
// @Param floor_id path string true "Floor ID"
// @Param zone_id path string true "Zone ID"
// @Description Get a zone from an area
// @Response 200 {object} GetZoneResponse

func GetZone(c *fiber.Ctx) error {
	areaInDB, err := db.GetArea(c)
	if err != nil {
		return err
	}
	floor, err := db.GetFloor(c, areaInDB)
	if err != nil {
		return err
	}
	zone, err := db.GetZone(c, floor)
	if err != nil {
		return err
	}
	return c.JSON(GetZoneResponse{Status: "success", Data: *zone})
}

type EditZoneResponse struct {
	Status string      `json:"status"`
	Data   models.Zone `json:"data"`
}

// EditZone edits a zone from a floor
// @Router /api/area/{area_id}/floor/{floor_id}/zone/{zone_id} [put]
// @Param area_id path string true "Area ID"
// @Param floor_id path string true "Floor ID"
// @Param zone_id path string true "Zone ID"
// @Description Edit a zone from a floor
// @Response 200 {object} EditZoneResponse
// @Security ApiKeyAuth
func EditZone(c *fiber.Ctx) error {
	areaInDB, err := db.GetAreaForUser(c)
	if err != nil {
		return err
	}

	floorInDB, err := db.GetFloor(c, areaInDB)
	if err != nil {
		return err
	}
	zoneInDB, err := db.GetZone(c, floorInDB)
	if err != nil {
		return err
	}
	zone := *zoneInDB // this is just nessesary because of the arrays used in the struct
	if err := c.BodyParser(&zone); err != nil {
		return err
	}
	zoneID := c.AllParams()["zone_id"]
	zoneId, err := primitive.ObjectIDFromHex(zoneID)
	newZoneBson := db.ConvertStructToBsonM(&zone, "updateble", "floors.$[floor].zones.$[zone].")

	filter := bson.D{{"_id", areaInDB.ID}}
	newZoneBson["updated_at"] = time.Now()
	update := bson.M{"$set": newZoneBson}

	// I hope I will never have to debug this
	arrayFilters := options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{
				{Key: "floor._id", Value: floorInDB.ID},
			},
			bson.D{
				{Key: "zone._id", Value: zoneId},
			},
		},
	})

	res := db.AreasCollection.FindOneAndUpdate(c.Context(), filter, update, arrayFilters)

	if res.Err() != nil {
		return db.ErrorHandler(err)
	}

	// todo read again from database to get the updated data
	var b bson.M
	res.Decode(&b)
	return GetZone(c)
}

type ZoneIdData struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Data     int                `json:"data" bson:"data"`
	Capacity int                `json:"capacity" bson:"capacity"`
	Time     time.Time          `json:"time" bson:"time"`
}

type GetSensorDataResponse struct {
	Status string       `json:"status"`
	Data   []ZoneIdData `json:"data"`
}

// GetSensorData gets the sensor data from a zone
// @Router /api/area/{area_id}/latest [get]
// @Param area_id path string true "Area ID"
// @Description Get the sensor data from a zone
// @Response 200 {object} GetSensorDataResponse
func GetSensorData(c *fiber.Ctx) error {
	areaInDB, err := db.GetArea(c)
	if err != nil {
		print("area not found")
		return err
	}
	var data []ZoneIdData
	for _, floor := range areaInDB.Floors {
		for _, zone := range floor.Zones {
			// query sensor in database that has the same id as the zone
			filter := bson.D{{"zone", zone.ID}}
			var sensor models.SensorModel
			err := db.SensorCollection.FindOne(c.Context(), filter).Decode(&sensor)
			if err != nil {
				continue
			}
			// get the last sensor data
			var sensorData models.SensorData
			//convert daystring to time
			filter = bson.D{{"sensor_id", sensor.ID}}
			from, err := time.Parse("2006-01-02T15:04:05.000Z", c.Query("from"))
			to, err2 := time.Parse("2006-01-02T15:04:05.000Z", c.Query("to"))
			if err == nil && err2 == nil {
				filter = append(filter, bson.E{Key: "time", Value: bson.D{{"$lte", to}, {"$gte", from}}})
			}
			opt := options.Find().SetSort(bson.D{{"time", -1}}).SetLimit(1)
			courser, err := db.SensorDataCollection.Find(c.Context(), filter, opt)
			found := false
			for courser.Next(c.Context()) {
				err := courser.Decode(&sensorData)
				if err != nil {
					return db.ErrorHandler(err)
				}
				found = true
			}
			if found {
				data = append(data, ZoneIdData{ID: zone.ID, Data: sensorData.Data, Time: sensorData.Time, Capacity: zone.Capacity})
			}

		}
	}
	return c.JSON(GetSensorDataResponse{Status: "success", Data: data})
}
