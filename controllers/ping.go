package controllers

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"rest-go/db"
)

func Ping(c *fiber.Ctx) error {
	const errMsg = "in /ping request the following error happened: "
	// testing mongodb connection
	err := db.DB.Client().Ping(c.Context(), nil)
	if err != nil {
		log.Fatalln(errMsg, err.Error())
	}
	return c.SendString("pong")
}
