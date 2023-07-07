package handlers

import (
	"chadgpt-api/app"
	"github.com/uptrace/bunrouter"
	"net/http"
)

type PlanHandler struct {
	app *app.App
}

type CreatePlanRequest struct {
	Diet string `json:"diet"`
}

func NewPlanHandler(app *app.App) *PlanHandler {
	return &PlanHandler{
		app: app,
	}
}

func (h *PlanHandler) Create(w http.ResponseWriter, req bunrouter.Request) error {
	client := h.app.GptClient()
	res, err := client.Request("/chat/completions")
	if err != nil {
		return err
	}

	return bunrouter.JSON(w, res)
}
