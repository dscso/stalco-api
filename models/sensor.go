package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SensorModel struct {
	ID  primitive.ObjectID `json:"id" bson:"_id"`
	Key string             `json:"key" bson:"key"`
}

type SensorData struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	SensorID primitive.ObjectID `json:"sensor_id" bson:"sensor_id"`
	Area     primitive.ObjectID `json:"area_id" bson:"area_id"`
}
