package mocks

import (
	"real-estate-system/listing-service/models"

	"github.com/stretchr/testify/mock"
)

type ListingRepositoryMock struct {
	mock.Mock
}

func (m *ListingRepositoryMock) CreateListing(listing *models.Listing) error {
	args := m.Called(listing)
	return args.Error(0)
}

func (m *ListingRepositoryMock) GetListings(page, size int) ([]models.Listing, error) {
	args := m.Called(page, size)
	return args.Get(0).([]models.Listing), args.Error(1)
}
