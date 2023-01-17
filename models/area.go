package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// todo logic for zone sensor relationship
type Zone struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name" api:"updateble"`
	Type     string             `json:"type" bson:"type" api:"updateble"`
	Radius   float32            `json:"radius" bson:"radius" api:"updateble"`
	Points   []float32          `json:"points" bson:"points" api:"updateble"`
	Capacity int                `json:"capacity" bson:"capacity" api:"updateble"`
}

type Floor struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Name   string             `json:"name" bson:"name" api:"updateble"`
	Number int                `json:"number" bson:"number" api:"updateble"`
	Image  string             `json:"image" bson:"image" api:"updateble"`
	Zones  []Zone             `json:"zones" bson:"zones"`
}

type Area struct {
	ID             primitive.ObjectID   `json:"id" bson:"_id"`
	Image          string               `json:"image" bson:"image" api:"updateble"`
	Description    string               `json:"description" bson:"description" api:"updateble"`
	Name           string               `json:"name" bson:"name" api:"updateble"`
	CreatedAt      time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at" bson:"updated_at"`
	Administrators []primitive.ObjectID `json:"-" bson:"administrators" api:"updateble"`
	Floors         []Floor              `json:"floors" bson:"floors"`
}
