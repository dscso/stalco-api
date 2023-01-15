package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SensorModel struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Name string             `json:"name" bson:"name" api:"updateble"`
	Type string             `json:"type" bson:"type" api:"updateble"`
	Key  string             `json:"-" bson:"key"`
	Area primitive.ObjectID `json:"area" bson:"area"`
}

type SensorData struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	SensorID primitive.ObjectID `json:"sensor_id" bson:"sensor_id"`
	Area     primitive.ObjectID `json:"area_id" bson:"area_id"`
}
