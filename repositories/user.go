package repositories

import (
	"github.com/MegeKaplan/megebase-identity-service/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (models.User, error)
	Create(user *models.User) error
}

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
