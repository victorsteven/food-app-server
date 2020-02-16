package interfaces

import (
	"fmt"
	"food-app/application"
	"food-app/domain/entity"
	"food-app/utils/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignIn(user *entity.User) (map[string]interface{}, map[string]string){

	var tokenErr = map[string]string{}
	//check if the user details are correct:
	u, err := application.UserApp.GetUserByEmailAndPassword(user)
	if err != nil {
		return nil, err
	}
	ts, tErr := token.CreateToken(u.ID)
	if tErr != nil {
		tokenErr["token_error"] = tErr.Error()
		return nil, err
	}
	saveErr := token.Auth.CreateAuth(u.ID, ts)
	if saveErr != nil {
		return nil, err
	}
	userData := make(map[string]interface{})
	userData["access_token"] = ts.AccessToken
	userData["refresh_token"] = ts.RefreshToken
	userData["id"] = u.ID
	userData["first_name"] = u.FirstName
	userData["last_name"] = u.LastName

	return userData, nil
}

func Login(c *gin.Context) {
	var user *entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	//validate request:
	validateUser := user.Validate("login")
	if len(validateUser) > 0 {
		c.JSON(http.StatusUnprocessableEntity, validateUser)
		return
	}
	userData, err := SignIn(user)
	if err != nil {
		//fmt.Println("the ")
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, userData)
}

func Logout(c *gin.Context) {

	//check is the user is authenticated first
	userId, err := token.Token.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized here")
		return
	}
	//userId, err := token.TokenAuth.FetchAuth(metadata.TokenUuid)
	//if err != nil {
	//	c.JSON(http.StatusUnauthorized, "unauthorized")
	//	return
	//}
	fmt.Print(userId)





	//var user *entity.User
	//if err := c.ShouldBindJSON(&user); err != nil {
	//	c.JSON(http.StatusUnprocessableEntity, "Invalid json")
	//	return
	//}
	//validate request:
	//validateUser := user.Validate("login")
	//if len(validateUser) > 0 {
	//	c.JSON(http.StatusUnprocessableEntity, validateUser)
	//	return
	//}
	//userData, err := SignIn(user)
	//if err != nil {
	//	//fmt.Println("the ")
	//	c.JSON(http.StatusInternalServerError, err)
	//	return
	//}
	c.JSON(http.StatusOK, "userData")
}
