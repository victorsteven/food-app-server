package interfaces

func (server *Server) Route()  {

	//router.POST("/login", interfaces.Login)

	server.Router.POST("/users", server.SaveUser)
	//router.GET("/users", interfaces.GetUsers)
	//router.GET("/users/:user_id", interfaces.GetUser)

	//If the file size is greater than 8MB dont allow it to even load into memory and waste our time.
	//router.POST("/food", middleware.AuthMiddleware(), middleware.MaxSizeAllowed(8192000), interfaces.SaveFood)
	//router.GET("/food", interfaces.GetAllFood)
	//router.GET("/food/:food_id", interfaces.GetFoodAndCreator)
	//router.PUT("/food/:food_id", middleware.AuthMiddleware(), middleware.MaxSizeAllowed(8192000), interfaces.UpdateFood)
	//router.DELETE("/food/:food_id", interfaces.DeleteFood)
}
