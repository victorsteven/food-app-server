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

func TestGetFood_Success(t *testing.T) {
	conn, err := Database()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	food, err := seedFood(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := NewFoodRepository(conn)

	f, saveErr := repo.GetFood(food.ID)

	assert.Nil(t, saveErr)
	assert.EqualValues(t, f.Title, food.Title)
	assert.EqualValues(t, f.Description, food.Description)
	assert.EqualValues(t, f.UserID, food.UserID)
}

func TestGetAllFood_Success(t *testing.T) {
	conn, err := Database()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	_, err = seedFoods(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := NewFoodRepository(conn)
	foods, getErr := repo.GetAllFood()

	assert.Nil(t, getErr)
	assert.EqualValues(t, len(foods), 2)
}

func TestUpdateFood_Success(t *testing.T) {
	conn, err := Database()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	food, err := seedFood(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	//updating
	food.Title = "food title update"
	food.Description = "food description update"

	repo := NewFoodRepository(conn)
	f, updateErr := repo.UpdateFood(food)

	assert.Nil(t, updateErr)
	assert.EqualValues(t, f.ID, 1)
	assert.EqualValues(t, f.Title, "food title update")
	assert.EqualValues(t, f.Description, "food description update")
	assert.EqualValues(t, f.UserID, 1)
}

//Duplicate title error
func TestUpdateFood_Failure(t *testing.T) {
	conn, err := Database()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	foods, err := seedFoods(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	var secondFood entity.Food

	//get the second food title
	for _, v := range foods {
		if v.ID == 1 {
			continue
		}
		secondFood = v
	}
	secondFood.Title = "first food" //this title belongs to the first food already, so the second food cannot use it
	secondFood.Description = "New description"

	repo := NewFoodRepository(conn)
	f, updateErr := repo.UpdateFood(&secondFood)

	dbMsg := map[string]string{
		"unique_title": "title already taken",
	}
	assert.NotNil(t, updateErr)
	assert.Nil(t, f)
	assert.EqualValues(t, dbMsg, updateErr)
}

func TestDeleteFood_Success(t *testing.T) {
	conn, err := Database()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	food, err := seedFood(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := NewFoodRepository(conn)

	deleteErr := repo.DeleteFood(food.ID)

	assert.Nil(t, deleteErr)
}