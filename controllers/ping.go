package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"rest-go/db"
	"rest-go/middleware"
	"rest-go/models"
)

func Ping(c *gin.Context) {
	const errMsg = "in /ping request the following error happened: "
	// testing mongodb connection
	err := db.DB.Client().Ping(c, nil)
	if err != nil {
		log.Fatalln(errMsg, err.Error())
	}
	// testing jwt
	user := models.User{
		Email: "hello@world.io",
	}
	token, err := middleware.GenerateJWT(user)
	if err != nil {
		log.Fatalln(errMsg, err.Error())
	}

	var userClaim models.UserClaim
	err = middleware.ValidateToken(token, &userClaim)
	if err != nil {
		log.Fatalln(errMsg, err.Error())
	}
	if user.Email != userClaim.Email {
		log.Fatalln(errMsg, "token is invalid")
	}
	c.String(http.StatusOK, "pong")
}
