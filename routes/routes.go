package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"

	"rest-go/controllers"
	_ "rest-go/docs"
	"rest-go/middleware"
)

//	@title			REST API for stalco Pils
//	@version		1.0
//	@description	backend for SPA Stalco
//	@termsOfService
//	@contact.name
//	@contact.email
//	@license.name
//	@license.url
//	@host
//	@BasePath	/api/

func StartServer() {

	app := fiber.New(
		fiber.Config{ErrorHandler: controllers.ErrorHandler},
	)
	middleware.InitSessions()

	app.Use(logger.New())

	// Default config
	app.Use(cors.New())

	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost, https://stalko.tk, http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	api := app.Group("/api")
	api.Use(middleware.Session)
	{
		api.Post("/user/signup", controllers.CreateUser)
		api.Post("/user/login", controllers.LoginUser)
		api.Get("/protected", controllers.Protected)

		api.Get("/area/:area_id", controllers.GetArea)
		api.Put("/area/:area_id", controllers.EditArea)

		api.Put("/area/:area_id/floors/:floor_id", controllers.EditFloor)
		api.Get("/area/:area_id/floors/:floor_id", controllers.GetFloor)

		api.Put("/area/:area_id/floors/:floor_id/zones/:zone_id", controllers.EditZone)
		api.Get("/area/:area_id/floors/:floor_id/zones/:zone_id", controllers.GetZone)

		api.Post("/area/:area_id/sensors", controllers.CreateSensor)
	}

	app.Get("/ping", controllers.Ping)

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	err := app.Listen(":8000")
	if err != nil {
		panic(err)
	}
}
