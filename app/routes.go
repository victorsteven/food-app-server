package app

import (
	"food-app/interfaces"
	"github.com/gin-gonic/gin"
)

func Route() *gin.Engine {

	router.POST("/users", interfaces.SaveUser)
	router.GET("/users", interfaces.GetUsers)
	router.GET("/users/:user_id", interfaces.GetUser)

	return router
}