package cmd

import (
	"food-app/cmd/middleware"
	"food-app/interfaces"
	"github.com/gin-gonic/gin"
)

func Route() *gin.Engine {

	router.POST("/login", interfaces.Login)

	router.POST("/users", interfaces.SaveUser)
	router.GET("/users", interfaces.GetUsers)
	router.GET("/users/:user_id", interfaces.GetUser)

	router.POST("/food", middleware.AuthMiddleware(), interfaces.SaveFood)
	router.GET("/food", interfaces.GetAllFood)
	//router.GET("/food/:food_id", interfaces.GetUser)

	return router
}
