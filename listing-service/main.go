package main

import (
	"fmt"
	"log"
	"os"
	"real-estate-system/listing-service/handlers"
	"real-estate-system/listing-service/models"
	"real-estate-system/listing-service/repository"
	"real-estate-system/listing-service/repository/interfaces"
	"real-estate-system/listing-service/seeders"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := buildDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	if err := db.AutoMigrate(&models.Listing{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	var repo interfaces.ListingRepository = repository.NewGormListingRepository(db)

	seeders.SeedListings(db)

	e := echo.New()
	handler := handlers.NewListingHandler(repo)

	e.GET("/listings", handler.GetListings)
	e.POST("/listings", handler.CreateListing)

	fmt.Println("Listing service running on :6000")
	e.Logger.Fatal(e.Start(":6000"))
}

func buildDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
		os.Getenv("DB_TIMEZONE"),
	)
}
