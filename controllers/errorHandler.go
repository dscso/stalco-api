package controllers

import "github.com/gofiber/fiber/v2"

func ErrorHandler(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a fiber.*Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Send error back to client
	return c.Status(code).JSON(fiber.Map{
		"status": "error",
		"error":  err.Error(),
	})
}
