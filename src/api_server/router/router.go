package router

import (
	"fmt"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/tylergeery/trash_hunt/api_server/controllers"
	"github.com/tylergeery/trash_hunt/api_server/middleware"
)

func health(c *routing.Context) error {
	return c.Write(fmt.Sprintf("Hello %s", c.Param("name")))
}

// GetRouter returns a new Mux Router
func GetRouter() *routing.Router {
	router := routing.New()

	router.Use(middleware.LogRequest(), middleware.Cors())
	router.Get(`/hello/<name:\w+>`, health)
	router.Post(
		"/login/",
		content.TypeNegotiator(content.JSON),
		controllers.PlayerLogin,
	)
	router.Post(
		"/v1/player/",
		content.TypeNegotiator(content.JSON),
		controllers.PlayerCreate,
	)

	rg := router.Group("/v1")
	rg.Use(
		middleware.ValidateToken(),
		content.TypeNegotiator(content.JSON),
	)

	rg.Put("/player/", controllers.PlayerUpdate)
	rg.Delete("/player/", controllers.PlayerDelete)
	rg.Get(`/player/<id:\d+>`, controllers.PlayerQuery)

	rg.Post("/game/start", controllers.GameStart)
	rg.Get("/game/status", controllers.GameStatus)
	rg.Post("/game/complete", controllers.GameComplete)

	rg.Get("/results/leaderboard/", controllers.Results)
	rg.Get("/results/me", controllers.MyResults)

	rg.Post("/auth", controllers.CreateAuthToken)

	return router
}
