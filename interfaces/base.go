package interfaces

import (
	"food-app/database/rdbms"
	"food-app/domain/entity"
	"food-app/utils/middleware"
	"food-app/utils/token"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
}

//

var (
	router = gin.Default()
)

type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

func (server *Server) DBConn() {

	//Connecting to a RDBMS
	dbdriver := os.Getenv("DB_DRIVER")
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	var err error

	server.DB, err = rdbms.NewDBConnection(dbdriver, user, password, port, host, dbname)
	if err != nil {
		log.Fatal("cannot connect to the db: ", err)
		return
	}
	//server.DB, err = infrastructure.NewDBConnection(dbdriver, user, password, port, host, dbname)
	//if err != nil {
	//	log.Fatal("cannot connect to the db: ", err)
	//	return
	//}
	server.DB.Debug().AutoMigrate(
		entity.User{},
		entity.Food{},
	)

	//Connecting to Redis
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASSWORD")

	_, err = token.TokenAuth.NewRedisClient(redis_host, redis_port, redis_password)
	if err != nil {
		log.Fatal(err)
	}

	server.Router = gin.Default()
	server.Router.Use(middleware.CORSMiddleware())

	server.Route()

	app_port := os.Getenv("PORT") //using heroku host
	if app_port == "" {
		app_port = "8888" //localhost
	}
	//log.Fatal(router.Run(":"+app_port))

}
func (server *Server) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
