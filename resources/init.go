package resources

import (
	"chadgpt-api/app"
	"chadgpt-api/resources/handlers"
	"context"
)

func init() {
	app.OnStart("resources.init", func(ctx context.Context, app *app.App) error {
		api := app.ApiRouter()
		userHandler := handlers.NewUserHandler(app)

		api.GET("/users", userHandler.List)
		api.GET("/users/:id", userHandler.Get)

		return nil
	})
}
