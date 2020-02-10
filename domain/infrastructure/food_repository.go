package infrastructure

import (
	"errors"
	"fmt"
	"food-app/domain/entity"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
	paginator "github.com/pilagod/gorm-cursor-paginator"
)

type repositoryFoodCRUD struct {
	db *gorm.DB
}
//const (
//	page = 2
//)

func NewRepositoryFoodCRUD(db *gorm.DB) *repositoryFoodCRUD {
	return &repositoryFoodCRUD{db}
}

func (r *repositoryFoodCRUD) SaveFood(food *entity.Food) (*entity.Food, error) {
	err := r.db.Debug().Create(&food).Error
	if err != nil {
		return nil, err
	}
	return food, nil
}

func (r *repositoryFoodCRUD) GetFood(id uint64) (*entity.Food, error) {
	var food entity.Food
	err := r.db.Debug().Where("id = ?", id).Take(&food).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		fmt.Println("error 2: ", err)
		return nil, errors.New("food not found")
	}

	return &food, nil
}

type PagingQuery struct {
	After  *string
	Before *string
	Limit  *int
	Order  *string
}

func GetModelPaginator(q PagingQuery) *paginator.Paginator {
	p := paginator.New()

	p.SetKeys( "Title") // [default: "ID"] (supporting multiple keys, order of keys matters)

	if q.After != nil {
		p.SetAfterCursor(*q.After) // [default: nil]
	}
	if q.Before != nil {
		p.SetBeforeCursor(*q.Before) // [default: nil]
	}
	if q.Limit != nil {
		p.SetLimit(*q.Limit) // [default: 10]
	}
	if q.Order != nil && *q.Order == "asc" {
		p.SetOrder(paginator.ASC) // [default: paginator.DESC]
	}
	return p
}

type Cursor struct {
	After  *string `json:"after"`
	Before *string `json:"before"`
}

func (r *repositoryFoodCRUD) GetAllFood(q PagingQuery) ([]entity.Food, paginator.Cursor,  error) {
	var foods []entity.Food
	stmt := r.db.Debug().Find(&foods)

	lim := new(int) //create a pointer of int
	*lim = 1

	before := new(string)
	*before = "2"
	after := new(string)
	*after = "2"

	//var q = PagingQuery{
	//	After:  after,
	//	Before: before,
	//	Limit:  lim,
	//	Order:  nil, //order is descending
	//}

	// get paginator for Model
	p := GetModelPaginator(q)

	result := p.Paginate(stmt, &foods)

	if result.Error != nil {
		fmt.Println("this is the error: ", result.Error)
	}
	// get cursor for next iteration
	cursor := p.GetNextCursor()

	return foods, cursor, nil
}

//func (r *repositoryFoodCRUD) GetAllFood() ([]entity.Food, error) {
//	var foods []entity.Food
//	//err := r.db.Debug().Find(&foods).Error
//	//if err != nil {
//	//	return nil, err
//	//}
//	//pag := pagination.Paging(&pagination.Param{
//	//	DB:      r.db,
//	//	Page:    page,
//	//	Limit:   limit,
//	//	OrderBy: []string{"id desc"},
//	//}, &foods)
//
//
//	//q := r.db.Debug()
//	////q := db.Model(Post{}).Where("published_at > ?", time.Now())
//	//p := paginator.New(adapter.NewGORMAdapter(q), 3)
//	//p.SetPage(1)
//	//_, _ = p.NextPage()
//	//if err := p.Results(&foods); err != nil {
//	//	panic(err)
//	//}
//
//	//if gorm.IsRecordNotFoundError(err) {
//	//	return nil, errors.New("user not found")
//	//}
//	//if len(foods) > 0 {
//	//	for i, _ := range foods {
//	//		err := r.db.Debug().Model(&entity.User{}).Where("id = ?", foods[i].UserID).Take(&foods[i].User).Error
//	//		if err != nil {
//	//			return nil, err
//	//		}
//	//	}
//	//}
//	//return foods, nil
//}

func (r *repositoryFoodCRUD) UpdateFood(food *entity.Food) (*entity.Food, error) {
	err := r.db.Debug().Save(&food).Error
	if err != nil {
		return nil, err
	}
	return food, nil
}

//func (r *repositoryFoodCRUD) GetFood(foodId uint64) (*entity.Food, error) {
//	var allfood []entity.Food
//	err := r.db.Debug().Find(&allfood).Error
//	if err != nil {
//		return nil, err
//	}
//	if gorm.IsRecordNotFoundError(err) {
//		return nil, errors.New("user not found")
//	}
//	return allfood, nil
//}
