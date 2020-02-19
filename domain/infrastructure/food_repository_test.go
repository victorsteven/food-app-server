package infrastructure

import (
	"food-app/domain/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSaveFood_Success(t *testing.T) {
	conn, err := Database()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	var food = entity.Food{}
	food.Title = "food title"
	food.Description = "food description"
	food.UserID = 1

	repo := NewFoodRepository(conn)

	f, saveErr := repo.SaveFood(&food)
	assert.Nil(t, saveErr)
	assert.EqualValues(t, f.Title, "food title")
	assert.EqualValues(t, f.Description, "food description")
	assert.EqualValues(t, f.UserID, 1)
}

//Failure can be due to duplicate email, etc
//Here, we will attempt saving a food that is already saved
func TestSaveFood_Failure(t *testing.T) {
	conn, err := Database()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	//seed the food
	_, err = seedFood(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	var food = entity.Food{}
	food.Title = "food title"
	food.Description = "food desc"
	food.UserID = 1

	repo := NewFoodRepository(conn)
	f, saveErr := repo.SaveFood(&food)

	dbMsg := map[string]string{
		"unique_title": "food title already taken",
	}
	assert.Nil(t, f)
	assert.EqualValues(t, dbMsg, saveErr)
}
