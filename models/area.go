package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Zone struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Type     string             `json:"type" bson:"string"`
	Radius   float32            `json:"radius" bson:"radius"`
	Points   []float32          `json:"points" bson:"points"`
	Capacity int                `json:"capacity" bson:"capacity"`
}

type Floor struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Name   string             `json:"name" bson:"name"`
	Number int                `json:"number" bson:"number"`
	Image  string             `json:"image" bson:"image"`
	Zones  []Zone             `json:"zones" bson:"zones"`
}

type Area struct {
	ID             primitive.ObjectID   `json:"id" bson:"_id"`
	Name           string               `json:"name" bson:"name"`
	CreatedAt      time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at" bson:"updated_at"`
	Administrators []primitive.ObjectID `json:"administrators" bson:"administrators"`
	Floors         []Floor              `json:"floors" bson:"floors"`
}
