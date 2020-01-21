package interfaces

import (
	"food-app/server/application"
	"food-app/server/domain/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func saveUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	u, err := application.NewUserImpl().SaveUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, u)
}

func getUser(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	u, err := application.NewUserImpl().GetUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, u)
}