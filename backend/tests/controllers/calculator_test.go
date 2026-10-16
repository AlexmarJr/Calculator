package controllers_test

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"calculator/backend/api/controllers"
	"calculator/backend/api/models"
	"calculator/backend/api/views"
)

type stubService struct {
	calculate func(models.Calculation) (string, error)
	apply     func(models.State, models.Action) (models.State, error)
}

func (s stubService) Calculate(input models.Calculation) (string, error) {
	return s.calculate(input)
}

func (s stubService) Apply(state models.State, action models.Action) (models.State, error) {
	return s.apply(state, action)
}

func TestHandlerDelegatesAndHidesUnexpectedErrors(t *testing.T) {
	service := stubService{calculate: func(input models.Calculation) (string, error) {
		if input.Operation != "add" || input.Left != "2" || input.Right == nil || *input.Right != "3" {
			t.Fatalf("incorrect command: %+v", input)
		}
		return "", errors.New("private infrastructure detail")
	}}
	request := httptest.NewRequest("POST", "/api/calculate", strings.NewReader(`{"operation":"add","left":"2","right":"3"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	controllers.NewCalculatorController(service).Calculate(response, request)
	if response.Code != 500 || strings.Contains(response.Body.String(), "private") {
		t.Fatalf("unsafe server error: %s", response.Body.String())
	}
}

func TestActionDTOMapsToDomainCommand(t *testing.T) {
	state := initialCalculatorState()
	state.History = []views.HistoryEntry{{ID: 1, Expression: "6 × 7", Result: "42"}}
	state.NextID = 2
	action := views.CalculatorAction{Type: "restore", Value: json.RawMessage(`{"id":1,"expression":"6 × 7","result":"42"}`)}
	body, _ := json.Marshal(views.ActionRequest{State: &state, Action: &action})
	called := false
	service := stubService{apply: func(state models.State, action models.Action) (models.State, error) {
		called = true
		if action.Type != "restore" || action.HistoryID != 1 || len(state.History) != 1 {
			t.Fatalf("incorrect command: %+v", action)
		}
		state.Display = "42"
		return state, nil
	}}
	request := httptest.NewRequest("POST", "/api/calculator/action", strings.NewReader(string(body)))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	controllers.NewCalculatorController(service).Action(response, request)
	if !called || response.Code != 200 {
		t.Fatalf("request was not delegated: %s", response.Body.String())
	}
	var output struct {
		State views.CalculatorState `json:"state"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &output); err != nil {
		t.Fatal(err)
	}
	if output.State.Display != "42" {
		t.Fatal("application result was not returned")
	}
}
