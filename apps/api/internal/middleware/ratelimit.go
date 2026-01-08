package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/vo"
)

// RateLimiter handles rate limiting
type RateLimiter struct {
	redis    *redis.Client
	requests int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(redisClient *redis.Client, requests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		redis:    redisClient,
		requests: requests,
		window:   window,
	}
}

// Middleware returns the rate limiting middleware
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip if Redis is not available
		if rl.redis == nil {
			c.Next()
			return
		}

		// Use client IP as the key
		key := fmt.Sprintf("ratelimit:%s", c.ClientIP())

		ctx := context.Background()
		
		// Get current count
		count, err := rl.redis.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			// Redis error, continue without rate limiting
			c.Next()
			return
		}

		if count >= rl.requests {
			c.JSON(http.StatusTooManyRequests, vo.ErrorVO{
				Code:    "RATE_LIMIT_EXCEEDED",
				Message: fmt.Sprintf("Rate limit exceeded. Try again in %v", rl.window),
			})
			c.Abort()
			return
		}

		// Increment counter
		pipe := rl.redis.Pipeline()
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, rl.window)
		_, _ = pipe.Exec(ctx)

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.requests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", rl.requests-count-1))

		c.Next()
	}
}

// InMemoryRateLimiter is a simple in-memory rate limiter for dev mode
type InMemoryRateLimiter struct {
	counts   map[string]int
	requests int
}

// NewInMemoryRateLimiter creates a new in-memory rate limiter
func NewInMemoryRateLimiter(requests int) *InMemoryRateLimiter {
	return &InMemoryRateLimiter{
		counts:   make(map[string]int),
		requests: requests,
	}
}

// Middleware returns the rate limiting middleware
func (rl *InMemoryRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()

		count := rl.counts[key]
		if count >= rl.requests {
			c.JSON(http.StatusTooManyRequests, vo.ErrorVO{
				Code:    "RATE_LIMIT_EXCEEDED",
				Message: "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		rl.counts[key] = count + 1

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.requests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", rl.requests-count-1))

		c.Next()
	}
}
