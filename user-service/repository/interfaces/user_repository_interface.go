package repository

import "real-estate-system/user-service/models"

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUsers(page, size int) ([]models.User, error)
	GetUser(id int) (*models.User, error)
}
