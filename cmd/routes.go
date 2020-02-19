package cmd

import (
	"food-app/cmd/middleware"
	"food-app/interfaces"
	"github.com/gin-gonic/gin"
)

func Route() *gin.Engine {

	router.POST("/login", interfaces.Login)
	router.POST("/logout", interfaces.Logout)
	router.POST("/refresh", interfaces.Refresh)



	router.POST("/users", interfaces.SaveUser)
	router.GET("/users", interfaces.GetUsers)
	router.GET("/users/:user_id", interfaces.GetUser)


	router.POST("/food", middleware.AuthMiddleware(), middleware.MaxSizeAllowed(8192000), interfaces.SaveFood)
	router.GET("/food", interfaces.GetAllFood)
	router.GET("/food/:food_id", interfaces.GetFoodAndCreator)
	router.PUT("/food/:food_id", middleware.AuthMiddleware(), middleware.MaxSizeAllowed(8192000), interfaces.UpdateFood)
	router.DELETE("/food/:food_id", interfaces.DeleteFood)

	return router
}