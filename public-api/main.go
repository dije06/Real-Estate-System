package main

import (
	"os"
	"real-estate-system/public-api/handlers"
	"time"

	custommiddleware "real-estate-system/public-api/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr(),
	})

	e.Use(custommiddleware.NewRedisRateLimiter(rdb, 5, time.Minute))

	// Public APIs
	e.POST("/public-api/users", handlers.CreateUser)
	e.POST("/public-api/listings", handlers.CreateListing)
	e.GET("/public-api/listings", handlers.GetListings)

	e.Logger.Fatal(e.Start(":6002"))
}

func redisAddr() string {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}
	return host + ":" + port
}
