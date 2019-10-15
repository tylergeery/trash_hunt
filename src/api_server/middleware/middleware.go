package middleware

import (
	"github.com/go-ozzo/ozzo-routing"
)

// LogRequest logs client request
func LogRequest() func(c *routing.Context) error {
	return logRequest()
}

// ValidateToken validates auth token and logs client request
func ValidateToken() func(c *routing.Context) error {
	return validateToken
}

// Cors sets expected CORS headers for requests
func Cors() func(c *routing.Context) error {
	return func(c *routing.Context) error {
		headers := c.Response.Header()
		headers.Set("Access-Control-Allow-Origin", "*")
		headers.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")

		if c.Request.Method == "OPTIONS" {
			headers.Set("Access-Control-Allow-Headers", "Authorization,Content-Type,If-Modified-Since,Cache-Control")
			headers.Set("Access-Control-Max-Age", "86400")
			c.Abort()
		}

		return nil
	}
}
