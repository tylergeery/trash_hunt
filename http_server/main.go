package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tylergeery/trash_hunt/http_server/controllers"
)

func health(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Fprintf(w, "Hello %s", r.URL.Path)
}

func main() {
	http.HandleFunc("/hello/", health)

	http.HandleFunc("/game/start", controllers.GameStart)
	http.HandleFunc("/game/status", controllers.GameStatus)
	http.HandleFunc("/game/complete", controllers.GameComplete)

	http.HandleFunc("/results/leaderboard/", controllers.Results)
	http.HandleFunc("/results/me", controllers.MyResults)

	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
