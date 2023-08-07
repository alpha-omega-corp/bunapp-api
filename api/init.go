package api

import (
	"context"
	"github.com/alpha-omega-corp/bunapp-api/api/handlers"
	"github.com/alpha-omega-corp/bunapp-api/app"
)

func Bootstrap() {
	app.OnStart("api.init", func(ctx context.Context, app *app.App) error {
		api := app.ApiRouter()
		userHandler := handlers.NewUserHandler(app)

		api.GET("/validate", userHandler.Validate)

		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)

		workoutHandler := handlers.NewWorkoutHandler(app)
		workoutGroup := api.NewGroup("/workout").Use(app.AuthHandler)

		workoutGroup.POST("/", workoutHandler.Create)

		return nil
	})
}
