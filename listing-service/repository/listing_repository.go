package repository

import (
	"real-estate-system/listing-service/models"

	"gorm.io/gorm"
)

type GormListingRepository struct {
	DB *gorm.DB
}

func NewGormListingRepository(db *gorm.DB) *GormListingRepository {
	return &GormListingRepository{DB: db}
}

func (r *GormListingRepository) CreateListing(listing *models.Listing) error {
	return r.DB.Create(listing).Error
}

func (r *GormListingRepository) GetListings(page, size int) ([]models.Listing, error) {
	var listings []models.Listing
	err := r.DB.Order("created_at desc").Limit(size).Find(&listings).Error
	return listings, err
}
