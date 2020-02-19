package rdbms

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
)

var db *gorm.DB

//Using Postgres
func NewDBConnection(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) (*gorm.DB, error) {
	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	db, err = gorm.Open(Dbdriver, DBURL)
	if err != nil {
		return nil, err
	} else {
		log.Println("CONNECTED TO: ", Dbdriver)
	}
	return db, nil
}

//Using MYSQL
//func  NewDBConnection(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) (*gorm.DB, error)  {
//	var err error
//	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
//
//	db, err = gorm.Open(Dbdriver, DBURL)
//	if err != nil {
//		log.Fatal("This is the error connecting to the database:", err)
//	}
//	fmt.Printf("We are connected to the %s database", Dbdriver)
//
//	return db, nil
//}

func NewDB() *gorm.DB {
	return db
}
