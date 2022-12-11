package db

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func ConvertStructToBsonM(structInput interface{}, tag string, bsonPrefix string) bson.M {
	newBson := bson.M{}
	fields := reflect.ValueOf(structInput).Elem()
	for i := 0; i < fields.NumField(); i += 1 {
		typeField := fields.Type().Field(i)
		// if string contains `tag`, add it to the bson
		if strings.Contains(typeField.Tag.Get("api"), tag) {
			newBson[bsonPrefix+typeField.Tag.Get("bson")] = fields.Field(i).Interface()
		}
	}
	return newBson
}
