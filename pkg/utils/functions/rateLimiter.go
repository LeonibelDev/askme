package functions

import (
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
)

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func ErrorHandler(c *gin.Context, info ratelimit.Info) {
	c.AbortWithStatusJSON(429, gin.H{
		"error":   "Too Many Requests",
		"message": "You have exceeded the maximum number of requests allowed. Please try again later. in seconds: " + time.Until(info.ResetTime).String(),
		"status":  429,
	})
}

func RateLimiter() gin.HandlerFunc {
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Second,
		Limit: 50,
	})

	mw := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: ErrorHandler,
		KeyFunc:      keyFunc,
	})

	return mw
}
