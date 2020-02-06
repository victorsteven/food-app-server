package interfaces

import (
	"food-app/application"
	"food-app/domain/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SaveFood(c *gin.Context) {
	var food *entity.Food
	if err := c.ShouldBindJSON(&food); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	u, err := application.FoodApp().SaveFood(food)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, u)
}

func GetAllFood(c *gin.Context) {
	allfood, err := application.FoodApp().GetAllFood()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, allfood)
}