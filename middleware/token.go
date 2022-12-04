package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"math/rand"
	"rest-go/models"
)

var secretKey = []byte("GoLinuxCasd,asndhjkashdahdkjahskdjahsdkjhasdjkahsdkjahsdkjashdjkashdkjashdakjsdhaksjhdloudKey")

func GenerateJWT(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{},
		ID:               rand.Int(),
		Email:            user.Email,
	})

	signedString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return "", fmt.Errorf("error creating signed string: %v", err)
	}

	return signedString, nil
}

func ValidateToken(tokenString string, userClaim *models.UserClaim) error {
	println("tokenString: ", tokenString)
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
