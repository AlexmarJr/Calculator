package controllers

import (
	"net/http"

	"calculator/backend/api/models"
	"calculator/backend/api/views"
)

// CalculatorService is the application boundary consumed by this HTTP adapter.
type CalculatorService interface {
	Calculate(models.Calculation) (string, error)
	Apply(models.State, models.Action) (models.State, error)
}

type CalculatorController struct {
	service CalculatorService
}

func NewCalculatorController(service CalculatorService) *CalculatorController {
	return &CalculatorController{service: service}
}

func (h *CalculatorController) Calculate(w http.ResponseWriter, r *http.Request) {
	var input views.CalculationRequest
	if !decodeJSON(w, r, 4096, &input, "Invalid calculation request.") {
		return
	}
	result, err := h.service.Calculate(input.ToModel())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	views.Calculation(w, result)
}

func (h *CalculatorController) Action(w http.ResponseWriter, r *http.Request) {
	var input views.ActionRequest
	if !decodeJSON(w, r, 32768, &input, "Invalid calculator action request.") {
		return
	}
	if input.State == nil || input.Action == nil {
		views.Error(w, http.StatusBadRequest, "State and action are required.")
		return
	}
	action, err := input.Action.ToModel()
	if err != nil {
		views.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	state, err := h.service.Apply(input.State.ToModel(), action)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	views.State(w, state)
}
