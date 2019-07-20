package main

import (
	"net/http"

	"github.com/tylergeery/trash_hunt/api_server/router"
)

func main() {
	router := router.GetRouter()

	http.Handle("/", router)
	panic(http.ListenAndServe(":8080", nil))
}
