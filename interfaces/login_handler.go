package interfaces

import (
	"fmt"
	"food-app/domain/entity"
	"food-app/utils/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)


type signinInterface interface {
	SignIn(*entity.User) (map[string]interface{}, map[string]string)
}
type sign struct {}

//var Sign signinInterface = &sign{} //The struct now implement the interface

//We will need to mock this method when writing unit test, it is best we define it in an interface.
//func (s *sign) SignIn(user *entity.User) (map[string]interface{}, map[string]string){
//	var tokenErr = map[string]string{}
//	//check if the user details are correct:
//	u, err := application.UserApp.GetUserByEmailAndPassword(user)
//	if err != nil {
//		return nil, err
//	}
//	ts, tErr := auth.Token.CreateToken(u.ID)
//	if tErr != nil {
//		tokenErr["token_error"] = tErr.Error()
//		return nil, err
//	}
//	saveErr := auth.Auth.CreateAuth(u.ID, ts)
//	if saveErr != nil {
//		return nil, err
//	}
//
//	userData := make(map[string]interface{})
//	userData["access_token"] = ts.AccessToken
//	userData["refresh_token"] = ts.RefreshToken
//	userData["id"] = u.ID
//	userData["first_name"] = u.FirstName
//	userData["last_name"] = u.LastName
//
//	return userData, nil
//}

//func Login(c *gin.Context) {
//	var user *entity.User
//	if err := c.ShouldBindJSON(&user); err != nil {
//		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
//		return
//	}
//	//validate request:
//	validateUser := user.Validate("login")
//	if len(validateUser) > 0 {
//		c.JSON(http.StatusUnprocessableEntity, validateUser)
//		return
//	}
//	userData, err := Sign.SignIn(user)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, err)
//		return
//	}
//	c.JSON(http.StatusOK, userData)
//}

func Logout(c *gin.Context) {
	//check is the user is authenticated first
	metadata, err := auth.Token.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	//if the access token exist and it is still valid, then delete both the access token and the refresh token
	deleteErr := auth.Auth.DeleteTokens(metadata)
	if deleteErr != nil {
		c.JSON(http.StatusUnauthorized, deleteErr.Error())
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}


//Refresh is the function that uses the refresh_token to generate new pairs of refresh and access tokens.
func Refresh(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	refreshToken := mapToken["refresh_token"]

	//verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//any error may be due to token expiration
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, "Cannot get uuid")
			return
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Error occurred")
			return
		}
		//Delete the previous Refresh Token
		 delErr := auth.Auth.DeleteRefresh(refreshUuid)
		if delErr != nil  { //if any goes wrong
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := auth.Token.CreateToken(userId)
		if  createErr != nil {
			c.JSON(http.StatusForbidden, createErr.Error())
			return
		}
		//save the tokens metadata to redis
		saveErr := auth.Auth.CreateAuth(userId, ts)
		if saveErr != nil {
			c.JSON(http.StatusForbidden, saveErr.Error())
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		c.JSON(http.StatusCreated, tokens)
	} else {
		c.JSON(http.StatusUnauthorized, "refresh token expired")
	}
}

