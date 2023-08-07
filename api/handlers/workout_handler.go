package handlers

import (
	"encoding/json"
	"github.com/alpha-omega-corp/bunapp-api/api/types"
	"github.com/alpha-omega-corp/bunapp-api/app"
	"github.com/uptrace/bunrouter"
	"net/http"
)

type WorkoutHandler struct {
	app *app.App
}

func NewWorkoutHandler(app *app.App) *WorkoutHandler {
	return &WorkoutHandler{
		app: app,
	}
}

func (h *WorkoutHandler) Create(w http.ResponseWriter, req bunrouter.Request) error {
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
