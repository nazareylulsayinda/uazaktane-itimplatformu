package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type client struct {
	limiter *rate.Limiter
}

var (
	mu      sync.Mutex
	clients = make(map[string]*client)
)

// RateLimit is an IP-based rate limiter middleware
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		mu.Lock()
		if _, found := clients[ip]; !found {
			// 5 istek/sn, burst 10
			clients[ip] = &client{limiter: rate.NewLimiter(5, 10)}
		}
		clientLimiter := clients[ip].limiter
		mu.Unlock()

		if !clientLimiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}
