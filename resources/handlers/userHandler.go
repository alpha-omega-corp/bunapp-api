package handlers

import (
	"chadgpt-api/app"
	"chadgpt-api/resources/models"
	"github.com/uptrace/bunrouter"
	"net/http"
)

type UserHandler struct {
	app *app.App
}

func NewUserHandler(app *app.App) *UserHandler {
	return &UserHandler{
		app: app,
	}
}

func (h *UserHandler) List(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()

	var users []models.User
	if err := h.app.Database().NewSelect().Model(&users).Scan(ctx); err != nil {
		return err
	}

	return bunrouter.JSON(w, bunrouter.H{
		"users": users,
	})
}

func (h *UserHandler) Get(w http.ResponseWriter, req bunrouter.Request) error {
	ctx := req.Context()
	id := req.Param("id")

	var user models.User
	if err := h.app.Database().NewSelect().Where("id = ?", id).Model(&user).Scan(ctx); err != nil {
		return err
	}

	return bunrouter.JSON(w, user)
}
