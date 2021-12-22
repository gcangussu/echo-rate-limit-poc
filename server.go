package main

import (
	"context"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","status":${status},"id":"${header:Id}","latency":${latency}}` + "\n",
	}))
	e.Use(RateLimiter())

	e.GET("/", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	e.Logger.Fatal(e.Start(":1323"))
}

var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

var limiter = redis_rate.NewLimiter(rdb)

func init() {
	rdb.FlushDB(context.TODO())
}

func RateLimiter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			idHeader := c.Request().Header["Id"]
			identifier := "___"
			if len(idHeader) > 0 {
				identifier = idHeader[0]
			}

			res, err := limiter.Allow(
				c.Request().Context(), identifier, redis_rate.PerSecond(10),
			)
			if err != nil {
				c.Logger().Errorf("Error on redis_rate limiter: %v", err)
				c.Error(c.NoContent(http.StatusInternalServerError))
				return nil
			}

			if res.Allowed <= 0 {
				c.Error(c.NoContent(http.StatusTooManyRequests))
				return nil
			}
			return next(c)
		}
	}
}
