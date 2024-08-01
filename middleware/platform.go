package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

func PlatformMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgentHeader := c.GetHeader("User-Agent")
		log.Printf("User-Agent: %s", userAgentHeader)
		userAgentHeader = strings.ToLower(userAgentHeader)

		if strings.Contains(userAgentHeader, "android") {
			c.Set("platform", "android")
		} else if strings.Contains(userAgentHeader, "iphone") || strings.Contains(userAgentHeader, "ios") {
			c.Set("platform", "ios")
		}

		c.Next()
	}
}
