package token

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"strconv"
	"time"
)

type tokenData struct {
	conn *redis.Client
}
type authInterface interface {
	CreateAuth(uint64, *TokenDetails) error
	FetchAuth(string) (uint64, error)
	DeleteTokens(string) (int64, error)
	NewRedisClient(host, port, password string) (*redis.Client, error)
}

var Auth authInterface = &tokenData{}

func (tk *tokenData) NewRedisClient(host, port, password string) (*redis.Client, error) {
	tk.conn = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0, //use default DB
	})
	return tk.conn, nil
}

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
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	atCreated, err := tk.conn.Set(td.TokenUuid, strconv.Itoa(int(userid)), at.Sub(now)).Result()
	if err != nil {
		return err
	}
	rtCreated, err := tk.conn.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Result()
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
	fmt.Println("WE ENTERED HERE")
	userid, err := tk.conn.Get(tokenUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

//Once a user row in the token table
func (tk *tokenData) DeleteTokens(tokenUuid string) (int64, error) {
	deleted, err := tk.conn.Del(tokenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
