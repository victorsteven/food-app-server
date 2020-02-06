package app

import (
	"fmt"
	"food-app/database"
	"food-app/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
}

var (
	router = gin.Default()
)


func StartApp() {
	//infrastructure.UserRepo.NewDBConnection(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	//route()
	dbdriver := os.Getenv("DB_DRIVER")
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	conn, err := database.NewDBConnection(dbdriver, user, password, port, host, dbname)
	if err != nil {
		log.Fatal("cannot connect to the db: ", err)
	} else {
		fmt.Println("DB CONNECTED")
	}
	conn.Debug().AutoMigrate(
		entity.User{},
	)

	Route()

	_ = router.Run(":8080")
}