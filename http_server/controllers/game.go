package controllers

import (
	"fmt"
	"net/http"
)

// GameStart - Create a new game
func GameStart(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Fprintf(w, "Hello %s", r.URL.Path)
}

// GameStatus - Read game status
func GameStatus(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Fprintf(w, "Hello %s", r.URL.Path)
}

// GameComplete - Finish a game and create results
func GameComplete(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Fprintf(w, "Hello %s", r.URL.Path)
}
