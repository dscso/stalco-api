package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"rest-go/middleware"
	"rest-go/util"

	"rest-go/db"
	"rest-go/models"
)

type CreateUserResponse struct {
	Status string      `json:"status" default:"success"`
	Data   models.User `json:"data"`
}

// CreateUser Creates a user
// @Router /api/user/signup [post]
// @Param User body models.User true "User"
// @Description Sign up as new user
// @Response 200 {object} CreateUserResponse
func CreateUser(c *fiber.Ctx) error {
	// converting json to struct
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return err
	}
	// todo better validation
	if len(user.Email) == 0 {
		return &fiber.Error{Message: "Email is required", Code: fiber.StatusBadRequest}
	}
	// hashing password
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error with password hash function: ", err.Error())
		return util.InternalServerError
	}
	user.Password = string(bytes)
	user.ID = primitive.NewObjectID()

	// inserting user in database
	_, err = db.UsersCollection.InsertOne(c.Context(), user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return &fiber.Error{Message: "User already exists", Code: fiber.StatusConflict}
		} else {
			log.Println("Mongo Error while inserting user: ", err.Error())
			return util.InternalServerError
		}
	}

	return c.JSON(CreateUserResponse{Status: "success", Data: user})
}

type LoginUserResponse struct {
	Status string                     `json:"status" default:"success"`
	Data   middleware.SessionDatabase `json:"data"`
}

// LoginUser logs in a user and returns a token
// @Router /api/user/login [post]
// @Param User body models.User true "User"
// @Description log in
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
		return db.ErrorHandler(err)
	}
	// comparing passwords
	if bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(user.Password)) != nil {
		return &fiber.Error{Message: "Incorrect password", Code: fiber.StatusUnauthorized}
	}

	sessionDataset, err := middleware.NewSession(c, userFromDB)
	if err != nil {
		return err
	}

	return c.JSON(LoginUserResponse{Status: "success", Data: *sessionDataset})
}

type UserObjectForResponse struct {
	ID    primitive.ObjectID `json:"id" bson:"_id"`
	Email string             `json:"email"`
	Areas []string           `json:"areas"`
}
type UserInfoResponse struct {
	Status  string                 `json:"status" default:"success"`
	Message string                 `json:"message"`
	User    *UserObjectForResponse `json:"user"`
}

func Protected(c *fiber.Ctx) error {
	user := c.Locals("session").(*util.SessionAuthenticated)

	if !user.Authenticated {
		return c.JSON(UserInfoResponse{Status: "error", Message: "Not authenticated"})
	}

	// query all areas, user is administrator of
	var areas []string
	cursor, err := db.AreasCollection.Find(c.Context(), bson.M{"administrators": user.UserID})
	if err != nil {
		return db.ErrorHandler(err)
	}
	for cursor.Next(c.Context()) {
		var area models.Area
		err := cursor.Decode(&area)
		if err != nil {
			return db.ErrorHandler(err)
		}
		areas = append(areas, area.ID.Hex())
	}
	// get email
	var userFromDB models.User
	err = db.UsersCollection.FindOne(c.Context(), bson.M{"_id": user.UserID}).Decode(&userFromDB)
	if err != nil {
		return db.ErrorHandler(err)
	}

	return c.JSON(UserInfoResponse{Status: "success", Message: "Authenticated", User: &UserObjectForResponse{ID: user.UserID, Email: userFromDB.Email, Areas: areas}})
}
