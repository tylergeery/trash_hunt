package controllers

import (
	"fmt"
	"net/http"
)

func Results(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Fprintf(w, "Hello %s", r.URL.Path)
}

func MyResults(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Fprintf(w, "Hello %s", r.URL.Path)
}
