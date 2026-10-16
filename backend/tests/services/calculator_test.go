package services_test

import (
	"errors"
	"testing"

	"calculator/backend/api/models"
	"calculator/backend/api/repositories"
	"calculator/backend/api/services"
)

func TestApplyValidatesIncomingState(t *testing.T) {
	_, err := services.NewCalculatorService(repositories.NewMemoryHistory).Apply(models.State{}, models.Action{Type: "clear"})
	var failure *models.Error
	if !errors.As(err, &failure) || failure.Kind != models.InvalidInput {
		t.Fatalf("invalid state accepted: %v", err)
	}
}

func TestApplyNormalizesHistoryAndReturnsDomainErrorsAsState(t *testing.T) {
	state := models.InitialState()
	state.History = nil
	next, err := services.NewCalculatorService(repositories.NewMemoryHistory).Apply(state, models.Action{Type: "clear"})
	if err != nil || next.History == nil {
		t.Fatalf("missing normalized history: %v", err)
	}
	state.Display = "-9"
	next, err = services.NewCalculatorService(repositories.NewMemoryHistory).Apply(state, models.Action{Type: "sqrt"})
	if err != nil || next.Error == nil || *next.Error != "Square root needs a non-negative number." {
		t.Fatalf("unexpected result: %+v, %v", next, err)
	}
}
