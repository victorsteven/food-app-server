package auth

import (
	"errors"
	"fmt"
	"food-app/database/redisdb"
	"strconv"
	"time"
)


type authInterface interface {
	CreateAuth(uint64, *TokenDetails) error
	FetchAuth(string) (uint64, error)
	DeleteRefresh(string) error
	DeleteTokens(*AccessDetails) error
}
type tokenData struct {}

var Auth authInterface = &tokenData{}

type AccessDetails struct {
	TokenUuid string
	UserId     uint64
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	TokenUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

//Save token metadata to Redis
func (tk *tokenData) CreateAuth(userid uint64, td *TokenDetails) error {
	conn := redisdb.NewRedisDB()
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	atCreated, err := conn.Set(td.TokenUuid, strconv.Itoa(int(userid)), at.Sub(now)).Result()
	if err != nil {
		return err
	}
	rtCreated, err := conn.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Result()
	if err != nil {
		return err
	}
	if atCreated == "0" || rtCreated == "0" {
		return errors.New("no record inserted")
	}
	return nil
}

//Check the metadata saved
func (tk *tokenData) FetchAuth(tokenUuid string) (uint64, error) {
	conn := redisdb.NewRedisDB()
	userid, err := conn.Get(tokenUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

//Once a user row in the token table
func (tk *tokenData) DeleteTokens(authD *AccessDetails) error {
	conn := redisdb.NewRedisDB()
	//get the refresh uuid
	refreshUuid := fmt.Sprintf("%s++%d", authD.TokenUuid, authD.UserId)
	//delete access token
	deletedAt, err := conn.Del(authD.TokenUuid).Result()
	if err != nil {
		return err
	}
	//delete refresh token
	deletedRt, err := conn.Del(refreshUuid).Result()
	if err != nil {
		return err
	}
	//When the record is deleted, the return value is 1
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("something went wrong")
	}
	return nil
}

func (tk *tokenData) DeleteRefresh(refreshUuid string) error {
	conn := redisdb.NewRedisDB()

	//delete refresh token
	deleted, err := conn.Del(refreshUuid).Result()
	if err != nil || deleted == 0 {
		return err
	}
	return nil
}

