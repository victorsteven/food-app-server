package cmd

import (
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

	dbdriver := os.Getenv("DB_DRIVER")
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	conn, err := database.NewDBConnection(dbdriver, user, password, port, host, dbname)
	if err != nil {
		log.Fatal("cannot connect to the db: ", err)
		return
	}
	conn.Debug().AutoMigrate(
		entity.User{},
		entity.Food{},
	)

	Route()

	_ = router.Run(":8080")
}