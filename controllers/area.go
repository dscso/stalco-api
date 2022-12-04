package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"rest-go/db"
	"rest-go/middleware"
	"rest-go/models"
)

func GetArea(c *gin.Context) {
	idString := c.Param("id")

	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		middleware.AppError(c, err, http.StatusNotFound, "Invalid ID")
		return
	}

	var ar models.Area

	err = db.AreasCollection.FindOne(c, bson.M{"_id": id}).Decode(&ar)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			middleware.AppError(c, err, http.StatusNotFound, "Area not found")
		} else {
			middleware.AppErrorFatal(c, err, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": ar})
}
