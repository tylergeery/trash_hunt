package controllers

import (
	"fmt"
	"net/http"
)

// Auth - Create a new auth token from user key
func Auth(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Fprintf(w, "Creating auth token %s", r.URL.Path)
}
