package middleware

import (
	"context"
	"log"
	"rest-go/util"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"rest-go/db"
	"rest-go/models"
)

var sessionTime = 60 * 60 * 24 * 7 // 1 week

type SessionAuthenticated struct {
	Authenticated bool
	UserID        primitive.ObjectID
}
type SessionDatabase struct {
	SessionId     primitive.ObjectID `json:"session_id" bson:"_id"`
	SessionSecret string             `json:"session_secret" bson:"session_secret"`
	UserID        primitive.ObjectID `json:"-" bson:"user_id"`
	CreatedAt     primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt     primitive.DateTime `json:"updated_at" bson:"updated_at"`
}

func InitSessions() {
	sparse := true
	expireAfter := int32(sessionTime)
	_, err := db.SessionsCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{"updated_at": 1},
		Options: &options.IndexOptions{
			Sparse:             &sparse,
			ExpireAfterSeconds: &expireAfter,
		},
	})
	if err != nil {
		panic(err)
	}
}

func Session(c *fiber.Ctx) error {
	sessAuth := SessionAuthenticated{Authenticated: false}
	c.Locals("session", &sessAuth)

	cookie := c.Cookies("auth", "None")
	if cookie == "None" {
		cookie = c.Get("Authorization", "None")
	}

	authString := strings.Split(cookie, " ")
	if len(authString) != 2 {
		return c.Next()
	}
	sessionId, err := primitive.ObjectIDFromHex(authString[0])
	if err != nil {
		// not authenticated
		return c.Next()
	}
	sessionSecret := authString[1]

	var sessionDataSet SessionDatabase
	if ok := db.SessionsCollection.FindOne(c.Context(), bson.M{"_id": sessionId}).Decode(&sessionDataSet); ok == nil {
		if sessionDataSet.SessionSecret == sessionSecret {
			// update last activity
			_, err = db.SessionsCollection.UpdateOne(c.Context(), bson.M{"_id": sessionId}, bson.M{"$set": bson.M{"updated_at": primitive.NewDateTimeFromTime(c.Context().Time())}})
			if err != nil {
				log.Println(err)
			}
			sessAuth.Authenticated = true
			sessAuth.UserID = sessionDataSet.UserID
		}
	}
	return c.Next()
}

func NewSession(c *fiber.Ctx, user models.User) (*SessionDatabase, error) {
	key, err := util.Create64ByteKey()
	if err != nil {
		return nil, err
	}

	sessionDataSet := SessionDatabase{
		SessionId:     primitive.NewObjectID(),
		SessionSecret: key,
		UserID:        user.ID,
		CreatedAt:     primitive.NewDateTimeFromTime(c.Context().Time()),
		UpdatedAt:     primitive.NewDateTimeFromTime(c.Context().Time()),
	}

	_, err = db.SessionsCollection.InsertOne(c.Context(), sessionDataSet)
	if err != nil {
		return nil, err
	}
	cookie := fiber.Cookie{
		Name:    "auth",
		Value:   sessionDataSet.SessionId.Hex() + " " + sessionDataSet.SessionSecret,
		Expires: c.Context().Time().Add(time.Duration(sessionTime) * time.Second),
	}

	c.Cookie(&cookie)
	return &sessionDataSet, nil
}
