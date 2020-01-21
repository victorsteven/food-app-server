package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

func NewDBConnection(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) *gorm.DB {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)

	db, err := gorm.Open(Dbdriver, DBURL)
	if err != nil {
		log.Fatalf("Cannot connect to the database: %v", err)
	}
	return db
}