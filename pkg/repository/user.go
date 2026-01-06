package repository

import (
	"github.com/La002/personal-crm/config"
	"github.com/La002/personal-crm/pkg/entity"
	"github.com/La002/personal-crm/pkg/logger"
	"github.com/La002/personal-crm/pkg/postgres"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(config *config.Configuration, l *logger.Logger) *UserRepo {
	db := postgres.ConnectDB(config, l)
	return &UserRepo{
		DB: db,
	}
}

func (r *UserRepo) CreateUser(user *entity.User) error {
	if err := r.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) GetUserByGoogleID(id string) (entity.User, error) {
	var user entity.User
	if err := r.DB.Where("google_id = ?", id).First(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *UserRepo) GetUserByEmail(email string) (entity.User, error) {
	var user entity.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *UserRepo) GetUserByID(id uint) (entity.User, error) {
	var user entity.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *UserRepo) UpdateUser(user *entity.User) error {
	if err := r.DB.Save(user).Error; err != nil {
		return err
	}
	return nil
}
