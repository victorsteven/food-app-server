package infrastructure

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
)
//var Mine
type Server struct {
	 DB *gorm.DB
}

//var DB *gorm.DB

//func (s *Server) NewDBConnection(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) (*gorm.DB, error) {
//	var err error
//	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
//	s.DB, err = gorm.Open(Dbdriver, DBURL)
//	if err != nil {
//		return nil, err
//	} else {
//		log.Println("CONNECTED TO: ", Dbdriver)
//	}
//	//return &Server{s.DB}, nil
//	return s.DB, nil
//}

//func  NewDBConnection(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) (*gorm.DB, error)  {
//	var err error
//	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
//
//	DB, err = gorm.Open(Dbdriver, DBURL)
//	if err != nil {
//		log.Fatal("This is the error connecting to the database:", err)
//	}
//	fmt.Printf("We are connected to the %s database", Dbdriver)
//
//	return DB, nil
//}

func  NewDBConnection(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) (*Server, error)  {
	var err error
	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)

	DB, err := gorm.Open(Dbdriver, DBURL)
	if err != nil {
		log.Fatal("This is the error connecting to the database:", err)
	}
	fmt.Printf("We are connected to the %s database", Dbdriver)

	return &Server{DB:DB}, nil
}

//func NewDB() *gorm.DB {
//	return DB
//}