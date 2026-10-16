package controllers_test

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"calculator/backend/api/views"
)

func testAction(t *testing.T, state views.CalculatorState, action views.CalculatorAction) views.CalculatorState {
	t.Helper()
	body, err := json.Marshal(views.ActionRequest{State: &state, Action: &action})
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest("POST", "/api/calculator/action", strings.NewReader(string(body)))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	newHandler().ServeHTTP(response, request)
	if response.Code != 200 {
		t.Fatalf("action %s returned %d: %s", action.Type, response.Code, response.Body.String())
	}
	var output struct {
		State views.CalculatorState `json:"state"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &output); err != nil {
		t.Fatal(err)
	}
	return output.State
}

func testValueAction(kind, value string) views.CalculatorAction {
	body, _ := json.Marshal(value)
	return views.CalculatorAction{Type: kind, Value: body}
}

func testFlow(t *testing.T, input string, state views.CalculatorState) views.CalculatorState {
	t.Helper()
	for _, key := range input {
		action := views.CalculatorAction{}
		switch {
		case key >= '0' && key <= '9':
			action = testValueAction("digit", string(key))
		case strings.ContainsRune("+-*/^", key):
			action = testValueAction("operator", string(key))
		default:
			kinds := map[rune]string{
				'=': "equals",
				'r': "sqrt",
				'%': "percent",
				'.': "decimal",
				',': "decimal",
				'S': "sign",
				'B': "backspace",
				'C': "clear",
				'H': "clear-history",
			}
			action.Type = kinds[key]
		}
		state = testAction(t, state, action)
	}
	return state
}

func TestServerCalculatorFlows(t *testing.T) {
	for _, item := range []struct {
		input, result string
	}{
		{"8+3=", "11"}, {"8-3=", "5"}, {"8*3=", "24"}, {"8/2=", "4"},
		{"2^10=", "1024"}, {"81r", "9"}, {"25%", "0.25"},
		{"200+10%=", "220"}, {"200-10%=", "180"}, {"200*10%=", "20"}, {"200/10%=", "2000"},
		{"0+10%=", "0"}, {"0.1+0.2=", "0.3"}, {"2+3*4=", "20"},
		{"2^^4=", "4"}, {"2^3^=", "3"}, {"2^*3=", "6"},
		{"2+9r4=", "6"}, {"8+S3=", "5"}, {"12.3B,4", "12.4"},
		{"1234.00", "1234.00"}, {"81r=", "9"}, {"8/0=7", "7"},
	} {
		t.Run(item.input, func(t *testing.T) {
			state := testFlow(t, item.input, initialCalculatorState())
			if state.Display != item.result || state.Error != nil {
				t.Fatalf("state: %+v", state)
			}
		})
	}
	for _, input := range []string{"8/0=", "9Sr", "2S^0.5="} {
		state := testFlow(t, input, initialCalculatorState())
		if state.Error == nil {
			t.Fatalf("missing domain error for %s", input)
		}
	}
}

func TestServerHistory(t *testing.T) {
	state := testFlow(t, "6*7=C", initialCalculatorState())
	if len(state.History) != 1 || state.Display != "0" {
		t.Fatal("clear should retain history")
	}
	entryJSON, _ := json.Marshal(state.History[0])
	state = testAction(t, state, views.CalculatorAction{Type: "restore", Value: entryJSON})
	state = testFlow(t, "+8=", state)
	if state.Display != "50" {
		t.Fatal("restore did not reuse the backend result")
	}
	state = testFlow(t, "H", state)
	if len(state.History) != 0 {
		t.Fatal("history was not cleared")
	}
	for i := 0; i < 25; i++ {
		state = testFlow(t, "C81r", state)
	}
	if len(state.History) != 20 {
		t.Fatal("history limit was not applied")
	}
	if err := validateState(state); err != nil {
		t.Fatal(err)
	}
}

func TestActionValidation(t *testing.T) {
	state := initialCalculatorState()
	for _, action := range []views.CalculatorAction{
		testValueAction("digit", "x"), testValueAction("operator", "%"),
		{Type: "unknown"}, {Type: "equals", Value: json.RawMessage(`"extra"`)},
		{Type: "restore", Value: json.RawMessage(`{"id":99,"expression":"1+1","result":"2"}`)},
	} {
		body, _ := json.Marshal(views.ActionRequest{State: &state, Action: &action})
		request := httptest.NewRequest("POST", "/api/calculator/action", strings.NewReader(string(body)))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		newHandler().ServeHTTP(response, request)
		if response.Code != 400 {
			t.Fatalf("invalid action returned %d", response.Code)
		}
	}
	for _, item := range []struct {
		body   string
		status int
	}{
		{`{}`, 400}, {`null`, 400},
		{`{"state":{},"action":{"type":"clear"}}`, 400},
		{strings.Repeat(" ", 33000), 413},
	} {
		request := httptest.NewRequest("POST", "/api/calculator/action", strings.NewReader(item.body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		newHandler().ServeHTTP(response, request)
		if response.Code != item.status {
			t.Fatalf("status %d, want %d", response.Code, item.status)
		}
	}
}

func TestConnectionErrorStartsFreshEntry(t *testing.T) {
	state := testFlow(t, "2+3", initialCalculatorState())
	state.Error = textPointer("Cannot reach the calculator service.")
	state = testAction(t, state, testValueAction("digit", "7"))
	if state.Display != "7" || state.Stored != nil || state.Operator != nil || state.Error != nil {
		t.Fatalf("error recovery retained pending calculation: %+v", state)
	}
	state = testFlow(t, "2+3", initialCalculatorState())
	state.Error = textPointer("Cannot reach the calculator service.")
	state = testAction(t, state, views.CalculatorAction{Type: "decimal"})
	if state.Display != "0." || state.Operator != nil {
		t.Fatal("decimal recovery retained pending calculation")
	}
}
