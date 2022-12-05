package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"rest-go/models"
	"strings"
	"time"
)

var secretKey = []byte("GoLinuxCasd,asndhjkashdahdkjahskdjahsdkjhasdjkahsdkjahsdkjashdjkashdkjashdakjsdhaksjhdloudKey")
var tokenDuration = time.Hour * 24

func GenerateJWT(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(tokenDuration)},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		},
		Email: user.Email,
	})

	signedString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return "", fmt.Errorf("error creating signed string: %v", err)
	}

	return signedString, nil
}

func ValidateToken(tokenString string, userClaim *models.UserClaim) error {
	token, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}

func AuthorizedJWT(c *gin.Context) {
	h := c.GetHeader("Authorization")
	split := strings.Split(h, " ")
	if len(split) != 2 {
		AppError(c, nil, http.StatusUnauthorized, "no token provided")
		return
	}
	var userClaim models.UserClaim
	err := ValidateToken(split[1], &userClaim)
	if err != nil {
		AppError(c, nil, http.StatusUnauthorized, "token invalid")
		return
	}

	c.Next()
}
