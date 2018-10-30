package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tylergeery/trash_hunt/http_server/router"
)

// GetControllerResponse performs http request and returns recorded result
func GetControllerResponse(t *testing.T, method, url string, body io.Reader, headers map[string]string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(method, url, body)
	router := router.GetRouter()

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	if err != nil {
		t.Fatalf("Error getting controller response: %s", err.Error())
	}

	router.ServeHTTP(rec, req)

	return rec
}
