package middleware

import (
	"net/http"
)

// LogRequest logs client request
func LogRequest(next http.Handler) http.Handler {
	return logRequest(next)
}

// LogRequestAndValidate validates auth token and logs client request
func LogRequestAndValidate(next http.Handler) http.Handler {
	return logRequest(validateToken(next))
}
