package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"rest-go/middleware"

	"rest-go/db"
	"rest-go/models"
)

type CreateUserResponse struct {
	Status string      `json:"status"`
	Data   models.User `json:"data"`
}

func CreateUser(c *fiber.Ctx) error {
	// converting json to struct
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return err
	}
	// hashing password
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	user.ID = primitive.NewObjectID()

	// inserting user in database
	_, err = db.UsersCollection.InsertOne(c.Context(), user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			println(err.Error())
			return &fiber.Error{Message: "User already exists", Code: fiber.StatusConflict}
		} else {
			return &fiber.Error{Message: "Internal server error", Code: fiber.StatusInternalServerError}
		}
	}

	return c.JSON(CreateUserResponse{Status: "success", Data: user})
}

type LoginUserResponse struct {
	Status string                     `json:"status"`
	Data   middleware.SessionDatabase `json:"data"`
}

func LoginUser(c *fiber.Ctx) error {
	// converting json to struct
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return err
	}
	// searching for user in database
	var userFromDB models.User
	err := db.UsersCollection.FindOne(c.Context(), bson.M{"email": user.Email}).Decode(&userFromDB)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &fiber.Error{Message: "User not found", Code: fiber.StatusNotFound}
		} else {
			return &fiber.Error{Message: "Internal server error", Code: fiber.StatusInternalServerError}
		}
	}
	// comparing passwords
	if bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(user.Password)) != nil {
		return &fiber.Error{Message: "Incorrect password", Code: fiber.StatusUnauthorized}
	}
	// saving user in session
	if err != nil {
		log.Println(err.Error())
		return &fiber.Error{Message: "Internal server error", Code: fiber.StatusInternalServerError}
	}
	sessionDataset, err := middleware.NewSession(c, userFromDB)
	if err != nil {
		return err
	}

	return c.JSON(LoginUserResponse{Status: "success", Data: *sessionDataset})
}

func Protected(c *fiber.Ctx) error {
	user := c.Locals("session").(*middleware.SessionAuthenticated)

	if user.Authenticated {
		return c.JSON(fiber.Map{"message": "Welcome " + user.UserID.Hex() + "!"})
	}
	return c.JSON(fiber.Map{"message": "Not authenticated!"})
}
