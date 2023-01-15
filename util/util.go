package util

import (
	"crypto/rand"
	"encoding/base64"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// there are no generics
func ContainsObjectID(ids []primitive.ObjectID, id primitive.ObjectID) bool {
	for _, i := range ids {
		if i == id {
			return true
		}
	}
	return false
}

func Create64ByteKey() (string, error) {
	key := make([]byte, 64)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

type SessionAuthenticated struct {
	Authenticated bool
	UserID        primitive.ObjectID
}
