package middleware

import (
	"net/http"
)

// LogRequestAndValidate validates auth token and logs client request
func LogRequestAndValidate(next http.Handler) http.Handler {
	return logRequest(validateToken(next))
}
