package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"rest-go/controllers"
)

// StartGin function
func StartGin() {
	router := gin.Default()
	api := router.Group("/api")
	{
		//api.GET("/area", user.GetAllUser)
		//api.POST("/users", user.CreateUser)
		api.GET("/area/:id", controllers.GetArea)
		/*api.PUT("/users/:id", user.UpdateUser)
		  api.DELETE("/users/:id", user.DeleteUser)*/
	}
	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})
	err := router.Run(":8000")
	if err != nil {
		panic(err)
	}
}
