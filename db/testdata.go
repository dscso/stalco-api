package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"rest-go/models"
)

func GenerateExampleArea() models.Area {
	return models.Area{
		ID:             primitive.NewObjectID(),
		Name:           "Marie Curie Bibliothek",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Administrators: []primitive.ObjectID{primitive.NewObjectID()},
		Floors: []models.Floor{
			models.Floor{
				ID:     primitive.NewObjectID(),
				Name:   "Ground Floor",
				Number: 0,
				Image:  "http://www.example.com/image.png",
				Zones: []models.Zone{
					models.Zone{
						ID:       primitive.NewObjectID(),
						Name:     "sensor one",
						Type:     "circle",
						Radius:   1.0,
						Capacity: 10,
					}, models.Zone{
						ID:       primitive.NewObjectID(),
						Name:     "Zone 2",
						Type:     "polygon",
						Points:   []float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0},
						Capacity: 100,
					},
				},
			}, models.Floor{
				ID:     primitive.NewObjectID(),
				Name:   "First Floor",
				Number: 1,
				Image:  "http://www.example.com/image.png",
				Zones: []models.Zone{
					models.Zone{
						ID:       primitive.NewObjectID(),
						Name:     "Zone 1",
						Type:     "circle",
						Radius:   1.0,
						Capacity: 10,
					}, models.Zone{
						ID:       primitive.NewObjectID(),
						Name:     "Zone 2",
						Type:     "polygon",
						Points:   []float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0},
						Capacity: 100,
					}, models.Zone{
						ID:       primitive.NewObjectID(),
						Name:     "Zone 3",
						Type:     "polygon",
						Points:   []float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0},
						Capacity: 100,
					},
				},
			},
		},
	}
}
