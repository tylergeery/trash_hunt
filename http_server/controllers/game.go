package controllers

import (
	"fmt"
	"net/http"
)

func GameStart(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Fprintf(w, "Hello %s", r.URL.Path)
}

func GameStatus(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Fprintf(w, "Hello %s", r.URL.Path)
}

func GameComplete(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Fprintf(w, "Hello %s", r.URL.Path)
}
