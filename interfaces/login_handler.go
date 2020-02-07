package interfaces

import (
	"fmt"
	"food-app/application"
	"food-app/domain/entity"
	"food-app/utils/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

func signin(user entity.User) (map[string]string, error) {

	//check if the user exist:
	u, err := application.UserApp().GetUserByEmailAndPassword(user.Email, user.Password)
	if err != nil {
		return nil, err
	}
	fmt.Println("THE ID TO SEND: ", u.ID)
	ts, err := token.CreateToken(u.ID)
	if err != nil {
		return nil, err
	}
	saveErr := token.TokenAuth.CreateAuth(u.ID, ts)
	if saveErr != nil {
		return nil, err
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	return tokens, nil
}

func Login(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	tokens, err := signin(user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.JSON(http.StatusOK, tokens)
}
