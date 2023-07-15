package resources

import (
	"chadgpt-api/app"
	"chadgpt-api/resources/handlers"
	"context"
)

func BootControllers() {
	app.OnStart("controller.init", func(ctx context.Context, app *app.App) error {
		api := app.ApiRouter()
		userHandler := handlers.NewUserHandler(app)

		api.GET("/user", app.UserFromToken)

		api.POST("/login", userHandler.Login)
		api.POST("/register", userHandler.Register)

		planHandler := handlers.NewPlanHandler(app)
		api.POST("/plans", app.AuthHandler(planHandler.Create))

		return nil
	})
}
