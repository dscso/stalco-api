package util

import "github.com/gofiber/fiber/v2"

var InternalServerError = &fiber.Error{Message: "Internal server error", Code: fiber.StatusInternalServerError}
var NotFoundError = &fiber.Error{Message: "Not found in the Database", Code: fiber.StatusNotFound}
var UnauthorizedError = &fiber.Error{Message: "Unauthorized", Code: fiber.StatusUnauthorized}
