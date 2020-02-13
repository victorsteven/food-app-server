package cmd

import (
	"food-app/interfaces"
	"github.com/joho/godotenv"
	"log"
)
//
func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
}

var server = interfaces.Server{}


func StartApp() {

	server.DBConn()

	//apiPort := fmt.Sprintf(":%s", os.Getenv("PORT"))
	//if apiPort == "" {
	//	apiPort = "8888"
	//}
	//fmt.Printf("Listening to port %s", apiPort)

	server.Run(":8888")
}
