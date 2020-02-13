package infrastructure

import (
	"fmt"
	"food-app/database/rdbms"
	"food-app/domain/entity"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	//Since we add our .env in .gitignore, Circle CI cannot see it, so see the else statement
	if _, err := os.Stat("./../../.env"); !os.IsNotExist(err) {
		var err error
		err = godotenv.Load(os.ExpandEnv("./../../.env"))
		if err != nil {
			log.Fatalf("Error getting env %v\n", err)
		}
		_, _ = Database()
	} else {
		//CIBuild()
		fmt.Println("nothing")
	}
	os.Exit(m.Run())
}

var server = repositoryUser{}

func Database() (*gorm.DB, error) {
	var err error
	dbdriver := os.Getenv("TEST_DB_DRIVER")
	host := os.Getenv("TEST_DB_HOST")
	password := os.Getenv("TEST_DB_PASSWORD")
	user := os.Getenv("TEST_DB_USER")
	dbname := os.Getenv("TEST_DB_NAME")
	port := os.Getenv("TEST_DB_PORT")

	server.db, err = rdbms.NewDBConnection(dbdriver, user, password, port, host, dbname)
	if err != nil {
		return nil, err
	}
	return server.db, nil
}

//Drop test db data if exist:
func refreshUserTable() error {
	err := server.db.DropTableIfExists(&entity.User{}, &entity.Food{}).Error
	if err != nil {
		return err
	}
	err = server.db.AutoMigrate(&entity.User{}, &entity.Food{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneUser() (*entity.User, error) {
	user := &entity.User{
		Email: "frank@gmail.com",
		FirstName: "victor",
		LastName: "steven",
		Password: "password",
	}
	err := server.db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}