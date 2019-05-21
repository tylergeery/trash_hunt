package main

import (
	"log"
	"net/http"

	"github.com/tylergeery/trash_hunt/api_server/router"
)

func main() {
	router := router.GetRouter()

	err := http.ListenAndServe(":8080", router) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
