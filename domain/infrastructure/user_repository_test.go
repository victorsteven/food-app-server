package infrastructure

import (
	"food-app/database/rdbms"
	"food-app/domain/entity"
	"github.com/jinzhu/gorm"
	"os"
)

func Database() (*gorm.DB, error) {
	dbdriver := os.Getenv("TEST_DB_DRIVER")
	host := os.Getenv("TEST_DB_HOST")
	password := os.Getenv("TEST_DB_PASSWORD")
	user := os.Getenv("TEST_DB_USER")
	dbname := os.Getenv("TEST_DB_NAME")
	port := os.Getenv("TEST_DB_PORT")

	conn, err := rdbms.NewDBConnection(dbdriver, user, password, port, host, dbname)
	if err != nil {
		return nil, err
	}
	err = conn.DropTableIfExists(&entity.User{}).Error
	if err != nil {
		return nil, err
	}
	err = conn.Debug().AutoMigrate(
		entity.User{},
	).Error
	if err != nil {
		return nil, err
	}
	return conn, nil
}

//func TestUserRepo_SaveUser(t *testing.T) {
//	conn, err := Database()
//	if err != nil {
//		t.Fatalf("want non error, got %#v", err)
//	}
//	var user = entity.User{}
//	user.ID = 1
//	user.Email = "manaan@gmail.com"
//	user.FirstName = "Kedu"
//	user.LastName = "Manner"
//	user.Password = "password"
//
//	repo := NewUserRepository(conn)
//
//	u, err := repo.SaveUser(&user)
//	if err != nil {
//		t.Fatalf("want non error, got %#v", err)
//	}
//}
