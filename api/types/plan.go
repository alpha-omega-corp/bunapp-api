package types

type CreatePlanRequest struct {
	Diet       string   `json:"diet"`
	Allergies  []string `json:"allergies"`
	Conditions []string `json:"conditions"`
	Goal       string   `json:"goal"`
	Bmi        float64  `json:"bmi"`
}
