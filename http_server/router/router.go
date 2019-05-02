package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tylergeery/trash_hunt/http_server/controllers"
	"github.com/tylergeery/trash_hunt/http_server/middleware"
)

func health(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Fprintf(w, "Hello %s", r.URL.Path)
}

func m(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(f)
}

// GetRouter returns a new Mux Router
func GetRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/hello/", health)

	router.HandleFunc("/login/", middleware.LogRequest(m(controllers.PlayerLogin)))

	router.Handle("/player/create", middleware.LogRequest(m(controllers.PlayerCreate)))
	router.Handle("/player/update", middleware.LogRequestAndValidate(m(controllers.PlayerUpdate)))
	router.Handle("/player/delete", middleware.LogRequestAndValidate(m(controllers.PlayerDelete)))
	router.Handle("/player/", middleware.LogRequestAndValidate(m(controllers.PlayerQuery)))

	router.Handle("/game/start", middleware.LogRequestAndValidate(m(controllers.GameStart)))
	router.Handle("/game/status", middleware.LogRequestAndValidate(m(controllers.GameStatus)))
	router.Handle("/game/complete", middleware.LogRequestAndValidate(m(controllers.GameComplete)))

	router.Handle("/results/leaderboard/", middleware.LogRequestAndValidate(m(controllers.Results)))
	router.Handle("/results/me", middleware.LogRequestAndValidate(m(controllers.MyResults)))

	router.Handle("/auth", middleware.LogRequestAndValidate(m(controllers.CreateAuthToken)))

	return router
}
