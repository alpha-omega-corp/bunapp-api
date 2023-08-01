package handler

import (
	"chadgpt-api/app"
	"context"
)

func Bootstrap() {
	app.OnStart("handler.init", func(ctx context.Context, app *app.App) error {
		api := app.ApiRouter()
		userHandler := NewUserHandler(app)

		api.POST("/login", userHandler.Login)
		api.POST("/register", userHandler.Register)

		api.GET("/user", userHandler.UserFromToken)

		planHandler := NewPlanHandler(app)
		api.POST("/plans", app.AuthHandler(planHandler.Create))

		return nil
	})
}