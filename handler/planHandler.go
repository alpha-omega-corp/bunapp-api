package handler

import (
	"encoding/json"
	"github.com/alpha-omega-corp/bunapp-api/app"
	"github.com/alpha-omega-corp/bunapp-api/types"
	"github.com/uptrace/bunrouter"
	"net/http"
)

type PlanHandler struct {
	app *app.App
}

func NewPlanHandler(app *app.App) *PlanHandler {
	return &PlanHandler{
		app: app,
	}
}

func (h *PlanHandler) Create(w http.ResponseWriter, req bunrouter.Request) error {
	var data types.CreatePlanRequest
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		return err
	}

	client := h.app.GptClient()
	promptManager := h.app.PromptManager()

	prompt, err := promptManager.Execute("head.prompt", data)
	if err != nil {
		return err
	}

	res, err := client.UserRequest(prompt)
	if err != nil {
		return err
	}

	return bunrouter.JSON(w, res)
}
