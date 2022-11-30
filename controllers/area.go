package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	
	"rest-go/db"
	"rest-go/models"
	"rest-go/util"
)

func GetArea(c *gin.Context) {
	idString := c.Param("id")

	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		util.HandleConvertIDError(c)
		return
	}

	var ar models.Area

	err = db.AreasCollection.FindOne(c, bson.M{"_id": id}).Decode(&ar)

	if err != nil {
		util.HandleMongoDecodeError(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": ar})
}
