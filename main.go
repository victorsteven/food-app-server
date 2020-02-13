package main

import (
	"food-app/interfaces"
)

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
func main() {
	StartApp()
}
