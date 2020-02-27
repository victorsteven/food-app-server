package interfaces

import (
	"fmt"
	"food-app/application"
	"food-app/domain/entity"
	"food-app/utils/auth"
	"food-app/utils/fileupload"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Food struct {
	foodApp    application.FoodAppInterface
	userApp    application.UserAppInterface
	fileUpload fileupload.UploadFileInterface
	tk         auth.TokenInterface
	rd         auth.AuthInterface
}

//Food constructor
func NewFood(fApp application.FoodAppInterface, uApp application.UserAppInterface, fd fileupload.UploadFileInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Food {
	return &Food{
		foodApp:    fApp,
		userApp:    uApp,
		fileUpload: fd,
		rd:         rd,
		tk:         tk,
	}
}

func (fo *Food) SaveFood(c *gin.Context) {
	//check is the user is authenticated first
	metadata, err := fo.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//lookup the metadata in redis:
	userId, err := fo.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//We we are using a frontend(vuejs), our errors need to have keys for easy checking, so we use a map to hold our errors
	var saveFoodError = make(map[string]string)

	title := c.PostForm("title")
	description := c.PostForm("description")
	if fmt.Sprintf("%T", title) != "string" || fmt.Sprintf("%T", description) != "string" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "Invalid json",
		})
		return
	}
	//We initialize a new food for the purpose of validating: in case the payload is empty or an invalid data type is used
	emptyFood := entity.Food{}
	emptyFood.Title = title
	emptyFood.Description = description
	saveFoodError = emptyFood.Validate("")
	if len(saveFoodError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, saveFoodError)
		return
	}
	file, err := c.FormFile("food_image")
	if err != nil {
		saveFoodError["invalid_file"] = "a valid file is required"
		c.JSON(http.StatusUnprocessableEntity, saveFoodError)
		return
	}
	//check if the user exist
	_, err = fo.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}
	uploadedFile, err := fo.fileUpload.UploadFile(file)
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
	savedFood, saveErr := fo.foodApp.SaveFood(&food)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr)
		return
	}
	c.JSON(http.StatusCreated, savedFood)
}

func (fo *Food) UpdateFood(c *gin.Context) {
	//Check if the user is authenticated first
	metadata, err := fo.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	//lookup the metadata in redis:
	userId, err := fo.rd.FetchAuth(metadata.TokenUuid)
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
	user, err := fo.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	//check if the food exist:
	food, err := fo.foodApp.GetFood(foodId)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	//if the user id doesnt match with the one we have, dont update. This is the case where an authenticated user tries to update someone else post using postman, curl, etc
	//if user.ID != food.UserID {
	//	c.JSON(http.StatusUnauthorized, "you are not the owner of this food")
	//	return
	//}
	//Since this is an update request,  a new image may or may not be given.
	// If not image is given, an error occurs. We know this that is why we ignored the error and instead check if the file is nil.
	// if not nil, we process the file by calling the "UploadFile" method.
	// if nil, we used the old one whose path is saved in the database
	file, _ := c.FormFile("food_image")
	if file != nil {
		food.FoodImage, err = fo.fileUpload.UploadFile(file)
		//since i am using Digital Ocean(DO) Spaces to save image, i am appending my DO url here. You can comment this line since you may be using Digital Ocean Spaces.
		food.FoodImage = os.Getenv("DO_SPACES_URL") + food.FoodImage
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
	updatedFood, dbUpdateErr := fo.foodApp.UpdateFood(food)
	if dbUpdateErr != nil {
		c.JSON(http.StatusInternalServerError, dbUpdateErr)
		return
	}
	c.JSON(http.StatusOK, updatedFood)
}

func (fo *Food) GetAllFood(c *gin.Context) {
	allfood, err := fo.foodApp.GetAllFood()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, allfood)
}

func (fo *Food) GetFoodAndCreator(c *gin.Context) {
	foodId, err := strconv.ParseUint(c.Param("food_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	food, err := fo.foodApp.GetFood(foodId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	user, err := fo.userApp.GetUser(food.UserID)
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

func (fo *Food) DeleteFood(c *gin.Context) {
	metadata, err := fo.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	foodId, err := strconv.ParseUint(c.Param("food_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	_, err = fo.userApp.GetUser(metadata.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	err = fo.foodApp.DeleteFood(foodId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "food deleted")
}
