package interfaces

import "real-estate-system/listing-service/models"

type ListingRepository interface {
	CreateListing(*models.Listing) error
	GetListings(page, size int) ([]models.Listing, error)
}
