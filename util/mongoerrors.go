package util

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleConvertIDError(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid ID"})
}

func HandleMongoDecodeError(err error, c *gin.Context) {
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "error": "Area not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Internal server error"})
			log.Println(err)
		}
		return
	}

}
