package services_test

import (
	"strings"
	"testing"

	"calculator/backend/api/models"
	"calculator/backend/api/repositories"
	"calculator/backend/api/services"
)

func runKeys(t *testing.T, state models.State, keys string) models.State {
	t.Helper()
	for _, key := range keys {
		action := models.Action{}
		switch {
		case key >= '0' && key <= '9':
			action = models.Action{Type: "digit", Value: string(key), HasValue: true}
		case strings.ContainsRune("+-*/^", key):
			action = models.Action{Type: "operator", Value: string(key), HasValue: true}
		default:
			action.Type = map[rune]string{
				'=': "equals",
				'%': "percent",
				'r': "sqrt",
				'.': "decimal",
				'C': "clear",
			}[key]
		}
		var err error
		state, err = services.NewCalculatorService(repositories.NewMemoryHistory).Apply(state, action)
		if err != nil {
			t.Fatal(err)
		}
		if err := services.ValidateState(state); err != nil {
			t.Fatal(err)
		}
	}
	return state
}

func TestStateRules(t *testing.T) {
	for _, item := range []struct {
		keys, expected string
	}{
		{"2+3*4=", "20"}, {"200+10%=", "220"}, {"200-10%=", "180"},
		{"200*10%=", "20"}, {"200/10%=", "2000"},
		{"2^^4=", "4"}, {"2^3^=", "3"}, {"2+9r4=", "6"}, {"81r=", "9"},
	} {
		state := runKeys(t, models.InitialState(), item.keys)
		if state.Error != nil || state.Display != item.expected {
			t.Fatalf("%s: %+v", item.keys, state)
		}
	}
}

func TestHistoryAndDomainErrorState(t *testing.T) {
	state := runKeys(t, models.InitialState(), "6*7=C")
	if len(state.History) != 1 || state.Display != "0" {
		t.Fatal("clear must preserve completed history")
	}
	state, err := services.NewCalculatorService(repositories.NewMemoryHistory).Apply(state, models.Action{Type: "restore", HistoryID: state.History[0].ID})
	if err != nil {
		t.Fatal(err)
	}
	state = runKeys(t, state, "+8=")
	if state.Display != "50" {
		t.Fatal("history result was not reused")
	}
	state = runKeys(t, state, "C9/0=")
	if state.Error == nil || state.Operator != nil {
		t.Fatal("undefined result must reset the pending operation")
	}
	state = runKeys(t, state, "7")
	if state.Display != "7" || state.Error != nil {
		t.Fatal("new input must recover from a domain error")
	}
}
