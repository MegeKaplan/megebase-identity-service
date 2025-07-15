package repositories

import (
	"github.com/MegeKaplan/megebase-identity-service/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (models.User, error)
	Create(user *models.User) error
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

// MOCK
type userMockRepository struct {
	users map[string]models.User
}

func NewUserMockRepository() UserRepository {
	return &userMockRepository{users: make(map[string]models.User)}
}

func (r *userMockRepository) FindByEmail(email string) (models.User, error) {
	if user, exists := r.users[email]; exists {
		return user, nil
	}
	return models.User{}, gorm.ErrRecordNotFound
}

func (r *userMockRepository) Create(user *models.User) error {
	r.users[user.Email] = *user
	return nil
}
