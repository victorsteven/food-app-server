package interfaces

import (
	"food-app/application"
	"food-app/domain/entity"
	"food-app/utils/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SaveFood(c *gin.Context) {
	var food *entity.Food
	if err := c.ShouldBindJSON(&food); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json")
		return
	}
	tokenAuth, err := token.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	userId, err := token.TokenAuth.FetchAuth(tokenAuth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	food.UserID = userId
	fo, err := application.FoodApp().SaveFood(food)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, fo)
}

func GetAllFood(c *gin.Context) {
	allfood, err := application.FoodApp().GetAllFood()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, allfood)
}
