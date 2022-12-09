package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"rest-go/controllers"
	"rest-go/middleware"
)

// StartGin function
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
		api.Get("/area/:id", controllers.GetArea)
		/*api.PUT("/users/:id", user.UpdateUser)
		  api.DELETE("/users/:id", user.DeleteUser)*/
	}

	app.Get("/ping", controllers.Ping)

	err := app.Listen(":8000")
	if err != nil {
		panic(err)
	}
}
