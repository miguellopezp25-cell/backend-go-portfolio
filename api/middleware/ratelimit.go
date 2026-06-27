package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	limiters map[string]*ipLimiter
	rate     rate.Limit
	burst    int
	cleanup  time.Duration
}

func NewRateLimiter(rps int, burst int) *RateLimiter {
	rl := &RateLimiter{
		limiters: make(map[string]*ipLimiter),
		rate:     rate.Limit(rps),
		burst:    burst,
		cleanup:  10 * time.Minute,
	}

	go rl.cleanupLoop()

	return rl
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	il, exists := rl.limiters[ip]
	if !exists {
		il = &ipLimiter{
			limiter:  rate.NewLimiter(rl.rate, rl.burst),
			lastSeen: time.Now(),
		}
		rl.limiters[ip] = il
		return il.limiter
	}

	il.lastSeen = time.Now()
	return il.limiter
}

func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, il := range rl.limiters {
			if time.Since(il.lastSeen) > rl.cleanup {
				delete(rl.limiters, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func RateLimit(rps int, burst int) gin.HandlerFunc {
	rl := NewRateLimiter(rps, burst)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			c.Header("Retry-After", "1")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}
