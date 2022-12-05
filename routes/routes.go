package routes

import (
	"net/http"

	"github.com/dscso/sessions"
	"github.com/dscso/sessions/mongo/mongodriver"
	"github.com/gin-gonic/gin"

	"rest-go/controllers"
	"rest-go/db"
	"rest-go/middleware"
)

var sessionTimeout = 60 * 60 * 24 * 7 // 1 week

// StartGin function
func StartGin() {
	router := gin.Default()

	sessionStore := mongodriver.NewStore(db.SessionsCollection, sessionTimeout, true, []byte("secret"))
	router.Use(sessions.Sessions("session_id", sessionStore))

	api := router.Group("/api")
	{
		//api.GET("/area", user.GetAllUser)
		api.POST("/user/signup", controllers.CreateUser)
		api.POST("/user/login", controllers.LoginUser)
		//api.POST("/user/refresh", controllers.RefreshToken)
		api.GET("/area/:id", controllers.GetArea)
		/*api.PUT("/users/:id", user.UpdateUser)
		  api.DELETE("/users/:id", user.DeleteUser)*/
	}

	router.GET("/ping", controllers.Ping)
	api.GET("/protected", middleware.Authorized, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "Gained access to protected resource"})
	})
	router.NoRoute(func(c *gin.Context) {
		middleware.AppError(c, nil, http.StatusNotFound, "API endpoint not found")
	})
	err := router.Run(":8000")
	if err != nil {
		panic(err)
	}
}
