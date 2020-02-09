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

	//If the file size is greater than 8MB dont allow it to even load into memory and waste our time.
	router.POST("/food", middleware.AuthMiddleware(), middleware.MaxSizeAllowed(8192000), interfaces.SaveFood)
	router.GET("/food", interfaces.GetAllFood)
	router.GET("/food/:food_id", interfaces.GetFood)

	//router.GET("/food/:food_id", interfaces.GetUser)

	return router
}
