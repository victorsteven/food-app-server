package interfaces

import (
	"fmt"
	"food-app/application"
	"food-app/domain/entity"
	"food-app/utils/fileupload"
	"food-app/utils/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SaveFood(c *gin.Context) {

	title := c.PostForm("title")
	description := c.PostForm("description")
	if fmt.Sprintf("%T", title) != "string" || fmt.Sprintf("%T", description) != "string" {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json")
	}
	file, err := c.FormFile("food_image")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid File")
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
	uploadedFile, err := fileupload.Uploader.UploadFile(file)
	if err != nil {
		fmt.Println("the error: ", err)
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	var food = entity.Food{}
	food.UserID = userId
	food.Title = title
	food.Description = description
	food.AvatarPath = uploadedFile
	fo, err := application.FoodApp().SaveFood(&food)
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
	c.JSON(http.StatusOK, allfood)
}
