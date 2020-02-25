package main

import (
	"food-app/domain/infrastructure"
	"food-app/interfaces"
	"food-app/utils/auth"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

const (
	Dbdriver = "postgres"
	DbHost     = "127.0.0.1"
	DbPort     = "5432"
	DbName   = "food-app"
	DbPassword = "password"
	DbUser     = "steven"

	redis_host = "127.0.0.1"
	redis_port = "6379"
	redis_password = ""
)

func init() {
	//To load our environmental variables.
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
}



func main() {
	services, err := infrastructure.NewServices(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.Automigrate()


	redisService, err := auth.NewRedisDB(redis_host, redis_port, redis_password)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	tk  := auth.NewToken()

	users := interfaces.NewUsers(services.User, redisService.Auth, tk)
	foods := interfaces.NewFood(services.Food, services.User, redisService.Auth, tk)


	r.POST("/login", users.Login)
	r.POST("/logout", users.Logout)
	r.POST("/refresh", users.Refresh)

	r.POST("/users", users.SaveUser)
	r.GET("/users", users.GetUsers)
	r.GET("/users/:user_id", users.GetUser)

	r.POST("/food", foods.SaveFood)
	r.POST("/food/:food_id", foods.UpdateFood)
	r.GET("/food/:food_id", foods.GetFoodAndCreator)
	r.DELETE("/food/:food_id", foods.DeleteFood)
	r.GET("/food", foods.GetAllFood)


	r.Run(":8888")

}

