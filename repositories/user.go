package repositories

import (
	"fmt"

	"github.com/MegeKaplan/megebase-identity-service/models"
	"github.com/MegeKaplan/megebase-identity-service/utils"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (models.User, error)
	Create(user *models.User) error
	FindByID(id string) (models.User, error)
	SearchUsers(params utils.QueryParams) ([]models.User, error)
}

// GORM
type userGormRepository struct {
	db *gorm.DB
}

func NewUserGormRepository(db *gorm.DB) UserRepository {
	return &userGormRepository{db: db}
}

func (r *userGormRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.First(&user, "email = ?", email).Error
	return user, err
}

func (r *userGormRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userGormRepository) FindByID(id string) (models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	return user, err
}

func (r *userGormRepository) SearchUsers(params utils.QueryParams) ([]models.User, error) {
	var users []models.User
	query := r.db.Model(&models.User{})

	for key, value := range params.Filters {
		if key == "limit" || key == "offset" || key == "sort" {
			continue
		}
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	if params.Sort != "" {
		query = query.Order(params.Sort)
	}

	query = query.Limit(params.Limit).Offset(params.Offset)

	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

// MOCK
// type userMockRepository struct {
// 	users map[string]models.User
// }

// func NewUserMockRepository() UserRepository {
// 	return &userMockRepository{users: make(map[string]models.User)}
// }

// func (r *userMockRepository) FindByEmail(email string) (models.User, error) {
// 	if user, exists := r.users[email]; exists {
// 		return user, nil
// 	}
// 	return models.User{}, gorm.ErrRecordNotFound
// }

// func (r *userMockRepository) Create(user *models.User) error {
// 	r.users[user.Email] = *user
// 	return nil
// }
