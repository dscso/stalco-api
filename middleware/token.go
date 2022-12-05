package middleware

import (
	"net/http"

	"github.com/dscso/sessions"
	"github.com/gin-gonic/gin"
)

func Authorized(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("userID")
	if user == nil {
		AppError(c, nil, http.StatusUnauthorized, "please sign in")
		return
	}

	c.Next()
}
