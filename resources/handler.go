package resources

import (
	"chadgpt-api/app"
	"chadgpt-api/resources/handlers"
	"context"
)

func Init() {
	app.OnStart("resources.init", func(ctx context.Context, app *app.App) error {
		api := app.ApiRouter()
		userHandler := handlers.NewUserHandler(app)

		api.POST("/login", userHandler.Login)
		api.POST("/users", userHandler.Create)

		api.GET("/users", app.AuthHandler(userHandler.List))
		api.GET("/users/:id", app.AuthClaimHandler(userHandler.Get))

		return nil
	})
}
