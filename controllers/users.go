package controllers

import (
	"log"
	"net/http"

	"github.com/dscso/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"rest-go/db"
	"rest-go/middleware"
	"rest-go/models"
)

func CreateUser(c *gin.Context) {
	// converting json to struct
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		middleware.AppError(c, err, http.StatusBadRequest, "Bad Request")
		return
	}
	// hashing password
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		middleware.AppErrorFatal(c, err, http.StatusInternalServerError, "Internal server error")
		return
	}
	user.Password = string(bytes)

	// inserting user in database
	_, err = db.UsersCollection.InsertOne(c, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			middleware.AppError(c, err, http.StatusConflict, "User already exists")
		} else {
			middleware.AppErrorFatal(c, err, http.StatusInternalServerError, "Internal server error")
		}
		return
	}
	// returning user
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{}})
}

// jwt authentication with token generation
func LoginUser(c *gin.Context) {
	// converting json to struct
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		middleware.AppError(c, err, http.StatusBadRequest, "Bad Request")
		return
	}
	// searching for user in database
	var userFromDB models.User
	err = db.UsersCollection.FindOne(c, bson.M{"email": user.Email}).Decode(&userFromDB)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middleware.AppError(c, err, http.StatusNotFound, "User not found")
		} else {
			middleware.AppErrorFatal(c, err, http.StatusInternalServerError, "Internal server error")
		}
		return
	}
	// comparing passwords
	err = bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(user.Password))
	if err != nil {
		middleware.AppError(c, err, http.StatusUnauthorized, "Unauthorized")
		return
	}
	// saving user in session
	session := sessions.Default(c)
	session.Set("userID", userFromDB.ID)
	err = session.Save()
	if err != nil {
		log.Println(err.Error())
		middleware.AppError(c, err, http.StatusInternalServerError, "Internal server error")
		return
	}
	// returning token
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{}})
}
