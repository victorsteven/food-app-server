package cmd

import (
	"food-app/cmd/middleware"
	"food-app/database/rdbms"
	"food-app/database/redisdb"
	"food-app/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

//
func init() {
	//To load our environmental variables.
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
}

var (
	router = gin.Default()
)

func StartApp() {

	//Connecting to a RDBMS
	dbdriver := os.Getenv("DB_DRIVER")
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	conn, err := rdbms.NewDBConnection(dbdriver, user, password, port, host, dbname)
	if err != nil {
		log.Fatal("cannot connect to the db: ", err)
		return
	}
	conn.Debug().AutoMigrate(
		entity.User{},
		entity.Food{},
	)

	//Connecting to Redis
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASSWORD")

	_, err = redisdb.NewRedisClient(redis_host, redis_port, redis_password)
	if err != nil {
		log.Fatal(err)
	}
	router.Use(middleware.CORSMiddleware())

	//Route()

	app_port := os.Getenv("PORT") //using heroku host
	if app_port == "" {
		app_port = "8888" //localhost
	}
	log.Fatal(router.Run(":" + app_port))
}
