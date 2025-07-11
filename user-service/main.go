package main

import (
	"fmt"
	"log"
	"os"
	"real-estate-system/user-service/handlers"
	"real-estate-system/user-service/models"
	"real-estate-system/user-service/repository"
	"real-estate-system/user-service/seeders"

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

	// Auto-migrate table
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	// Seed
	seeders.SeedUsers(db)

	e := echo.New()

	userRepo := repository.NewGormUserRepository(db)

	h := handlers.NewUserHandler(userRepo)

	// Routes
	e.GET("/users", h.GetUsers)
	e.GET("/users/:id", h.GetUser)
	e.POST("/users", h.CreateUser)

	fmt.Println("User service running on :6001")
	e.Logger.Fatal(e.Start(":6001"))
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
