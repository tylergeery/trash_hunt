package middleware

// LogRequestAndValidate validates auth token and logs client request
func LogRequestAndValidate(next http.Hander) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(validateToken(next))
	})
}
