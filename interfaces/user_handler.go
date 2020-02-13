package interfaces

import (
	"food-app/application"
	"food-app/domain/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func SaveUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	u, err := application.UserApp.SaveUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, u.PublicUser())
}

func GetUsers(c *gin.Context) {
	us := entity.Users{} //customize user
	var err error
	us, err = application.UserApp.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, us.PublicUsers())
}

func GetUser(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user, err := application.UserApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user.PublicUser())
}
