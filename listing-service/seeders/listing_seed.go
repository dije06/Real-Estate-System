package seeders

import (
	"fmt"
	"math/rand"
	"real-estate-system/listing-service/models"
	"time"

	"gorm.io/gorm"
)

var listingTypes = []string{"rent", "sale"}

func randomListingType() string {
	return listingTypes[rand.Intn(len(listingTypes))]
}

func randomPrice() int {
	if rand.Intn(2) == 0 {
		return rand.Intn(3000) + 2000 // harga dummy rent: 2000–4999
	}
	return rand.Intn(90000) + 10000 // harga dummy sale: 10000–99999
}

func randomUserID() int {
	// Dummy User
	return rand.Intn(10) + 1
}

func SeedListings(db *gorm.DB) {
	var count int64
	if err := db.Model(&models.Listing{}).Count(&count).Error; err != nil {
		fmt.Println("Failed to count listings:", err)
		return
	}

	// if count > 0 {
	// 	fmt.Println("Listings already seeded.")
	// 	return
	// }

	rand.Seed(time.Now().UnixNano())
	now := time.Now().UnixMicro()

	for i := 0; i < 10; i++ {
		listing := models.Listing{
			UserID:      randomUserID(),
			Price:       randomPrice(),
			ListingType: randomListingType(),
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		db.Create(&listing)
	}

	fmt.Println("Seeded 10 random listings.")
}
