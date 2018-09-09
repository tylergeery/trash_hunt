package controllers

import (
	"fmt"
	"net/http"
)

// Results - get results for game
func Results(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Fprintf(w, "Hello %s", r.URL.Path)
}

// MyResults - get results for my games
func MyResults(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Fprintf(w, "Hello %s", r.URL.Path)
}
