package router

import (
	"fmt"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/cors"
	"github.com/tylergeery/trash_hunt/src/api_server/controllers"
	"github.com/tylergeery/trash_hunt/src/api_server/middleware"
)

func health(c *routing.Context) error {
	return c.Write(fmt.Sprintf("Hello %s", c.Param("name")))
}

// GetRouter returns a new Mux Router
func GetRouter() *routing.Router {
	router := routing.New()

	router.Get(`/hello/<name:\w+>`, health)
	router.Post(
		"/login/",
		middleware.LogRequest(),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.Options{
			AllowOrigins: "*",
			AllowHeaders: "*",
			AllowMethods: "*",
		}),
		controllers.PlayerLogin,
	)
	router.Post(
		"/v1/player/create",
		middleware.LogRequest(),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.Options{
			AllowOrigins: "*",
			AllowHeaders: "*",
			AllowMethods: "*",
		}),
		controllers.PlayerCreate,
	)

	rg := router.Group("/v1")
	rg.Use(
		middleware.LogRequest(),
		middleware.ValidateToken(),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.Options{
			AllowOrigins: "*",
			AllowHeaders: "*",
			AllowMethods: "*",
		}),
	)

	rg.Post("/player/update", controllers.PlayerUpdate)
	rg.Post("/player/delete", controllers.PlayerDelete)
	rg.Get("/player/", controllers.PlayerQuery)

	rg.Post("/game/start", controllers.GameStart)
	rg.Get("/game/status", controllers.GameStatus)
	rg.Post("/game/complete", controllers.GameComplete)

	rg.Get("/results/leaderboard/", controllers.Results)
	rg.Get("/results/me", controllers.MyResults)

	rg.Post("/auth", controllers.CreateAuthToken)

	return router
}
