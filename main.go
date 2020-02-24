package main

import (
	"fmt"
	"food-app/domain/infrastructure"
	"food-app/interfaces"
	"github.com/gin-gonic/gin"
)

const (
	DbHost     = "127.0.0.1"
	DbPort     = "5432"
	DbName   = "food-app"
	DbPassword = "password"
	DbUser     = "steven"
)

func main() {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)

	services, err := infrastructure.NewServices(DBURL)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.Automigrate()

	//r := mux.NewRouter()
	r := gin.Default()

	usersC := interfaces.NewUsers(services.User)

	r.POST("/users", usersC.SaveUser)
	r.GET("/users", usersC.GetUsers)
	r.GET("/users/:user_id", usersC.GetUser)

	r.Run(":8888")

	//galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

	//isProd := false
	//b, err := rand.Bytes(32)
	//must(err)
	//csrfMw := csrf.Protect(b, csrf.Secure(isProd))
	//fmt.Println("this is the csrf: ", csrfMw)
	////http.ListenAndServe(":8000", CSRF(r))
	//userMw := middleware.User{
	//	UserService: services.User,
	//}
	////This is the user middleware
	//requireUserMw := middleware.RequireUser{User: userMw}
	//
	//r.Handle("/", staticC.Home).Methods("GET")
	//r.Handle("/contact", staticC.Contact).Methods("GET")
	//
	//r.HandleFunc("/signup", usersC.New).Methods("GET")
	//r.HandleFunc("/signup", usersC.Create).Methods("POST")
	//
	//r.Handle("/login", usersC.LoginView).Methods("GET")
	//r.HandleFunc("/login", usersC.Login).Methods("POST")
	////r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	//
	////Asset route
	////the dir path specified below is the path from where the app root resides
	//assertHandler := http.FileServer(http.Dir("./assets/"))
	//assertHandler = http.StripPrefix("/assets/", assertHandler)
	//r.PathPrefix("/assets/").Handler(assertHandler)
	//
	////Make the image appear on the screen
	//imageHandler := http.FileServer(http.Dir("./images/"))
	//r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))
	//
	////Gallery Routes
	////New is not a function. It is actually a view
	//r.Handle("/galleries/new", requireUserMw.Apply(galleriesC.New)).Methods("GET")
	//r.Handle("/galleries", requireUserMw.ApplyFn(galleriesC.Index)).Methods("GET")
	//
	////Create is a function.
	//r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
	//r.HandleFunc("/galleries/{id}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET").Name(controllers.EditGallery)
	//r.HandleFunc("/galleries/{id}/update", requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
	//r.HandleFunc("/galleries/{id}/delete", requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")
	//
	//r.HandleFunc("/galleries/{id}/images", requireUserMw.ApplyFn(galleriesC.ImageUpload)).Methods("POST")
	//r.HandleFunc("/galleries/{id}/images/{filename}/delete", requireUserMw.ApplyFn(galleriesC.ImageDelete)).Methods("POST")
	//
	//r.HandleFunc("/galleries/{id}", galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)
	//
	//fmt.Println("this is the one: ", csrfMw(userMw.Apply(r)))
	//if err := http.ListenAndServe(":3000", csrfMw(userMw.Apply(r))); err != nil {
	//	log.Fatal("cannot start server")
	//}
}

