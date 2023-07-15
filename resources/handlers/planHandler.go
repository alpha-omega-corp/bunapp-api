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
	promptManager := h.app.PromptManager()

	prompt, err := promptManager.Execute("head.prompt", struct {
		Diet       string
		Allergies  []string
		Conditions []string
		Goal       string
		Bmi        float64
	}{
		Diet:       "Keto",
		Allergies:  []string{"Peanuts", "Shellfish"},
		Conditions: []string{"Diabetes", "High Blood Pressure"},
		Goal:       "Lose Weight",
		Bmi:        25.0,
	})
	if err != nil {
		return err
	}

	res, err := client.UserRequest(prompt)
	if err != nil {
		return err
	}

	return bunrouter.JSON(w, res)
}
