package middleware

import (
	"time"

	"github.com/dekanayake/kart-challenge/backend-challenge/internal/config"
	"github.com/gin-gonic/gin"
)

func ZerologMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		config.Logger.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("queryParams", c.Request.URL.RawQuery).
			Int("status", status).
			Str("client_ip", c.ClientIP()).
			Dur("latency", latency).
			Msg("HTTP request")
	}
}
