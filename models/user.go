package models

import "github.com/golang-jwt/jwt/v4"

type User struct {
	ID       string `json:"id" bson:"_id"`
	Email    string `json:"email" bson:"email" required:"true"`
	Password string `json:"password" required:"true"`
}

type UserClaim struct {
	jwt.RegisteredClaims
	//ID    int    `json:"id"`
	Email string `json:"email"`
}
