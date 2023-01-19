package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SensorModel struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Name string             `json:"name" bson:"name" api:"updateble"`
	Type string             `json:"type" bson:"type" api:"updateble"`
	Key  string             `json:"-" bson:"key"`
	Area primitive.ObjectID `json:"area" bson:"area"`
	Zone primitive.ObjectID `json:"zone" bson:"zone" api:"updateble"`
}

type SensorData struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	SensorID primitive.ObjectID `json:"sensor_id" bson:"sensor_id"`
	Time     time.Time          `json:"time" bson:"time"`
	Data     int                `json:"data" bson:"data"`
}
