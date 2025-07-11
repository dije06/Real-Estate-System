package repository

import (
	"real-estate-system/user-service/models"

	"gorm.io/gorm"
)

type GormUserRepository struct {
	DB *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{DB: db}
}

func (r *GormUserRepository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *GormUserRepository) GetUsers(page, size int) ([]models.User, error) {
	var users []models.User
	offset := (page - 1) * size
	result := r.DB.Order("created_at desc").Offset(offset).Limit(size).Find(&users)
	return users, result.Error
}

func (r *GormUserRepository) GetUser(id int) (*models.User, error) {
	var user models.User
	result := r.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
