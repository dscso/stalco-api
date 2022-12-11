package util

import "go.mongodb.org/mongo-driver/bson/primitive"

// there are no generics
func ContainsObjectID(ids []primitive.ObjectID, id primitive.ObjectID) bool {
	for _, i := range ids {
		if i == id {
			return true
		}
	}
	return false
}
