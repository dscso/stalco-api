package routes

import (
	"github.com/gofiber/fiber/v2"
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
	api := app.Group("/api")
	api.Use(middleware.Session)

	{
		//api.GET("/area", user.GetAllUser)
		api.Post("/user/signup", controllers.CreateUser)
		api.Post("/user/login", controllers.LoginUser)
		api.Get("/protected", controllers.Protected)
		//api.POST("/user/refresh", controllers.RefreshToken)
		api.Get("/area/:area_id", controllers.GetArea)
		//api.Post("/area/:area_id/floors", controllers.EditFloorFactory(fiber.MethodPost))
		api.Put("/area/:area_id/floors/:floor_id", controllers.EditFloor)

		/*api.PUT("/users/:id", user.UpdateUser)
		  api.DELETE("/users/:id", user.DeleteUser)*/
	}

	app.Get("/ping", controllers.Ping)

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	err := app.Listen(":8000")
	if err != nil {
		panic(err)
	}
}
