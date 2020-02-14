package interfaces

import (
	"fmt"
	"food-app/application"
	"food-app/domain/entity"
	"food-app/utils/fileupload"
	"food-app/utils/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func SaveFood(c *gin.Context) {
	//check is the user is authenticated first
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
	//We we are using a frontend(vuejs), our errors need to have keys for easy checking, so we use a map to hold our errors
	var saveFoodError = make(map[string]string)

	title := c.PostForm("title")
	description := c.PostForm("description")
	if fmt.Sprintf("%T", title) != "string" || fmt.Sprintf("%T", description) != "string" {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json")
	}
	//We initialize a new food for the purpose of validating: in case the payload is empty or an invalid data type is used
	emptyFood := entity.Food{}
	emptyFood.Title = title
	emptyFood.Description = description
	saveFoodError = emptyFood.Validate("update")
	if len(saveFoodError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, saveFoodError)
		return
	}
	file, err := c.FormFile("food_image")
	if err != nil {
		saveFoodError["invalid_file"] = "invalid file"
		c.JSON(http.StatusUnprocessableEntity, saveFoodError)
		return
	}
	uploadedFile, err := fileupload.Uploader.UploadFile(file)
	if err != nil {
		saveFoodError["upload_err"] = err.Error() //this error can be any we defined in the UploadFile method
		c.JSON(http.StatusUnprocessableEntity, saveFoodError)
		return
	}
	var food = entity.Food{}
	food.UserID = userId
	food.Title = title
	food.Description = description
	food.FoodImage = uploadedFile
	fo, err := application.FoodApp.SaveFood(&food)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, fo)
}

func UpdateFood(c *gin.Context) {
	//Check if the user is authenticated first
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
	//We we are using a frontend(vuejs), our errors need to have keys for easy checking, so we use a map to hold our errors
	var updateFoodError = make(map[string]string)

	foodId, err := strconv.ParseUint(c.Param("food_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	//Since it is a multipart form data we sent, we will do a manual check on each item
	title := c.PostForm("title")
	description := c.PostForm("description")
	if fmt.Sprintf("%T", title) != "string" || fmt.Sprintf("%T", description) != "string" {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json")
	}
	//We initialize a new food for the purpose of validating: in case the payload is empty or an invalid data type is used
	emptyFood := entity.Food{}
	emptyFood.Title = title
	emptyFood.Description = description
	updateFoodError = emptyFood.Validate("update")
	if len(updateFoodError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, updateFoodError)
		return
	}
	//check if the food exist:
	food, err := application.FoodApp.GetFood(foodId)
	if err != nil {
		c.JSON(http.StatusNotFound, "food not found")
		return
	}
	//if the user id doesnt match with the one we have, dont update. This is the case where an authenticated user tries to update someone else post using postman, curl, etc
	if userId != food.UserID {
		c.JSON(http.StatusUnauthorized, "you are not the owner of this food")
		return
	}
	//Since this is an update request,  a new image may or may not be given.
	// If not image is given, an error occurs. We know this that is why we ignored the error and instead check if the file is nil.
	// if not nil, we process the file by calling the "UploadFile" method.
	// if nil, we used the old one whose path is saved in the database
	file, _ := c.FormFile("food_image")
	if file != nil {
		food.FoodImage, err = fileupload.Uploader.UploadFile(file)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"upload_err": err.Error(),
			})
			return
		}
	}
	//we dont need to update user's id
	food.Title = title
	food.Description = description
	food.UpdatedAt = time.Now()
	fo, err := application.FoodApp.UpdateFood(food)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"db_err": "too large: upload an image less than 8MB",
		})
		return
	}
	c.JSON(http.StatusOK, fo)
}

func GetAllFood(c *gin.Context) {
	allfood, err := application.FoodApp.GetAllFood()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, allfood)
}

func GetFoodAndCreator(c *gin.Context) {
	foodId, err := strconv.ParseUint(c.Param("food_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	food, err := application.FoodApp.GetFood(foodId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	user, err := application.UserApp.GetUser(food.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	foodAndUser := map[string]interface{}{
		"food":    food,
		"creator": user.PublicUser(),
	}
	c.JSON(http.StatusOK, foodAndUser)
}

func DeleteFood(c *gin.Context) {
	foodId, err := strconv.ParseUint(c.Param("food_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	err = application.FoodApp.DeleteFood(foodId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "food deleted")
}
