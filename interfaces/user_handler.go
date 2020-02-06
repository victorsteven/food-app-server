package interfaces

import (
	"fmt"
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
	u, err := application.UserApp().SaveUser(&user)
	if err != nil {
		fmt.Println("WE HAVE AN ERROR")
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, u)
}

func GetUsers(c *gin.Context) {
	users, err := application.UserApp().GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, users)
}

func  GetUser(c *gin.Context) {
	_, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//u, err := interactor.GetUser(userId)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, err.Error())
	//	return
	//}
	c.JSON(http.StatusCreated, "dfs")
}