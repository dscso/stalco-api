package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func AppError(c *gin.Context, err error, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{"status": "error", "error": message})
}

func AppErrorFatal(c *gin.Context, err error, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{"status": "error", "error": message})
	log.Fatal(err)
}
