package seeders

import (
	"fmt"
	"math/rand"
	"real-estate-system/user-service/models"
	"time"

	"gorm.io/gorm"
)

// Sample Indonesian first and last names
var FirstNames = []string{
	"Agus", "Dewi", "Rizky", "Indah", "Putra", "Siti", "Bayu", "Citra", "Eka", "Yuni",
	"Fajar", "Rahma", "Andi", "Lestari", "Rina", "Joko", "Nur", "Budi", "Arif", "Mega",
}

var LastNames = []string{
	"Pratama", "Wijaya", "Saputra", "Santoso", "Siregar", "Wulandari", "Setiawan", "Maulana",
	"Kusuma", "Hardian", "Ramadhan", "Febrianto", "Handayani", "Putri", "Utami", "Susanto",
}

func randomName() string {
	first := FirstNames[rand.Intn(len(FirstNames))]
	last := LastNames[rand.Intn(len(LastNames))]
	return fmt.Sprintf("%s %s", first, last)
}

func SeedUsers(db *gorm.DB) {
	var count int64

	// Count the number of existing records
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		fmt.Println("Failed to count users:", err)
		return
	}

	// Skip seeding if already populated
	if count > 0 {
		fmt.Println("Users already seeded.")
		return
	}

	// Proceed with seeding if empty
	fmt.Println("Seeding users...")
	rand.Seed(time.Now().UnixNano())
	now := time.Now().UnixMicro()

	for i := 0; i < 10; i++ {
		user := models.User{
			Name:      randomName(),
			CreatedAt: now,
			UpdatedAt: now,
		}
		db.Create(&user)
	}

	fmt.Println("Seeded 10 users.")
}
