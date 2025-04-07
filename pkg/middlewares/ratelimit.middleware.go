package middlewares

import (
	"sync"
	"time"

	"github.com/dinhdev-nu/realtime_auth_go/pkg/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(5, 10)

func RateLimitMiddleware() gin.HandlerFunc { // 5 request per second, burst 10
	return func(c *gin.Context) {
		if !limiter.Allow() {
			response.TooManyRequests(c)
			c.Abort()
			return
		}
		c.Next()
	}
}

// tự làm
func RateLimitMiddleware2() gin.HandlerFunc {
	type RateLimitPerSecond struct {
		mu       sync.Mutex // mutex để đảm bảo thread safety khi cập nhiều luồng  api cập nhật counter 1 lúc
		lastTime time.Time
		counter  int
	}

	var limiterRate = &RateLimitPerSecond{}

	return func(c *gin.Context) {
		now := time.Now()

		limiterRate.mu.Lock() // có thể sleep 1 chút để tránh deadlock

		if now.Sub(limiterRate.lastTime) > time.Second {
			limiterRate.lastTime = now
			limiterRate.counter = 0
		} else {
			limiterRate.counter++
			if limiterRate.counter > 10 {
				response.TooManyRequests(c)
				c.Abort()
				return

			}
		}

		limiterRate.mu.Unlock()

		c.Next()
	}
}
