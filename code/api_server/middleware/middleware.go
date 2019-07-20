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
