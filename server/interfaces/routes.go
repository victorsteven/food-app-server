package interfaces

import "github.com/gin-gonic/gin"

func Route() *gin.Engine {

	router := gin.Default()

	router.POST("/users", saveUser)
	router.GET("/users/:id", getUser)

	return router
}