package interfaces

import (
	"food-app/application"
	"food-app/domain/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)


//func NewUserApp(db *gorm.DB) {
//
//}


func (server *Server) SaveUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	//validate the request:
	validateErr := user.Validate("")
	if len(validateErr) > 0 {
		c.JSON(http.StatusUnprocessableEntity, validateErr)
		return
	}
	app := application.UserImpl{DB: server.DB}
	u, err := app.SaveUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, u.PublicUser())
}

//func GetUsers(c *gin.Context) {
//	us := entity.Users{} //customize user
//	var err error
//	us, err = application.UserApp().GetUsers()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, err.Error())
//		return
//	}
//	c.JSON(http.StatusCreated, us.PublicUsers())
//}
//
//func GetUser(c *gin.Context) {
//	userId, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, err.Error())
//		return
//	}
//	user, err := application.UserApp().GetUser(userId)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, err.Error())
//		return
//	}
//	c.JSON(http.StatusOK, user.PublicUser())
//}
