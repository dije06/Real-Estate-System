package mocks

import (
	"real-estate-system/user-service/models"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepositoryMock) GetUsers(page, size int) ([]models.User, error) {
	args := m.Called(page, size)
	var users []models.User
	if args.Get(0) != nil {
		users = args.Get(0).([]models.User)
	}
	return users, args.Error(1)
}

func (m *UserRepositoryMock) GetUser(id int) (*models.User, error) {
	args := m.Called(id)
	var user *models.User
	if args.Get(0) != nil {
		user = args.Get(0).(*models.User)
	}
	return user, args.Error(1)
}
