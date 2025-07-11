package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func NewRedisRateLimiter(rdb *redis.Client, limit int, window time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			key := fmt.Sprintf("ratelimit:%s", ip)
			ctx := context.Background()

			// Increment count and set expiration if it's a new key
			count, err := rdb.Incr(ctx, key).Result()
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Rate limiter error")
			}
			if count == 1 {
				rdb.Expire(ctx, key, window)
			}

			if int(count) > limit {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"message": "Rate limit exceeded",
				})
			}

			return next(c)
		}
	}
}
