package handler

import (
	"chadgpt-api/app"
	"context"
)

func Init() {
	app.OnStart("handler.init", func(ctx context.Context, app *app.App) error {
		api := app.ApiRouter()
		userHandler := NewUserHandler(app)

		api.GET("/user", app.UserFromToken)

		api.POST("/login", userHandler.Login)
		api.POST("/register", userHandler.Register)

		planHandler := NewPlanHandler(app)
		api.POST("/plans", app.AuthHandler(planHandler.Create))

		return nil
	})
}
