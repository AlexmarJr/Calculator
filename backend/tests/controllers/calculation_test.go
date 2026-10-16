package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"calculator/backend/api/views"
)

func pointer(value string) *string {
	return &value
}

func TestCalculationAPI(t *testing.T) {
	tests := []struct {
		name    string
		input   views.CalculationRequest
		result  string
		status  int
		message string
	}{
		{"add", views.CalculationRequest{Operation: "add", Left: "8", Right: pointer("3")}, "11", 200, ""},
		{"subtract", views.CalculationRequest{Operation: "subtract", Left: "8", Right: pointer("3")}, "5", 200, ""},
		{"multiply", views.CalculationRequest{Operation: "multiply", Left: "8", Right: pointer("3")}, "24", 200, ""},
		{"divide", views.CalculationRequest{Operation: "divide", Left: "8", Right: pointer("3")}, "2.66666666666667", 200, ""},
		{"power", views.CalculationRequest{Operation: "power", Left: "8", Right: pointer("3")}, "512", 200, ""},
		{"sqrt", views.CalculationRequest{Operation: "sqrt", Left: "81"}, "9", 200, ""},
		{"percent", views.CalculationRequest{Operation: "percent", Left: "25"}, "0.25", 200, ""},
		{"relative percent", views.CalculationRequest{Operation: "percent", Left: "10", Base: pointer("200")}, "20", 200, ""},
		{"zero base", views.CalculationRequest{Operation: "percent", Left: "10", Base: pointer("0")}, "0", 200, ""},
		{"decimal precision", views.CalculationRequest{Operation: "add", Left: "0.1", Right: pointer("0.2")}, "0.3", 200, ""},
		{"negative operand", views.CalculationRequest{Operation: "add", Left: "8", Right: pointer("-3")}, "5", 200, ""},
		{"unfinished decimal", views.CalculationRequest{Operation: "add", Left: "2.", Right: pointer("3")}, "5", 200, ""},
		{"negative zero", views.CalculationRequest{Operation: "multiply", Left: "-0", Right: pointer("3")}, "0", 200, ""},
		{"scientific operand", views.CalculationRequest{Operation: "add", Left: "1e6", Right: pointer("1")}, "1000001", 200, ""},
		{"small decimal", views.CalculationRequest{Operation: "divide", Left: "1", Right: pointer("100000")}, "0.00001", 200, ""},
		{"division zero", views.CalculationRequest{Operation: "divide", Left: "5", Right: pointer("0")}, "", 422, "Cannot divide by zero."},
		{"negative root", views.CalculationRequest{Operation: "sqrt", Left: "-9"}, "", 422, "Square root needs a non-negative number."},
		{"fractional negative power", views.CalculationRequest{Operation: "power", Left: "-2", Right: pointer("0.5")}, "", 422, "The result is not a real number."},
		{"overflow", views.CalculationRequest{Operation: "power", Left: "99", Right: pointer("999")}, "", 422, "The result is too large."},
		{"invalid operation", views.CalculationRequest{Operation: "unknown", Left: "2"}, "", 400, "Unsupported operation."},
		{"missing right", views.CalculationRequest{Operation: "add", Left: "2"}, "", 400, ""},
		{"unexpected right", views.CalculationRequest{Operation: "sqrt", Left: "2", Right: pointer("3")}, "", 400, ""},
		{"unexpected base", views.CalculationRequest{Operation: "add", Left: "2", Right: pointer("3"), Base: pointer("5")}, "", 400, ""},
		{"formatted number rejected", views.CalculationRequest{Operation: "sqrt", Left: "1.234,56"}, "", 400, ""},
		{"infinite operand", views.CalculationRequest{Operation: "sqrt", Left: "1e999"}, "", 400, ""},
		{"NaN rejected", views.CalculationRequest{Operation: "sqrt", Left: "NaN"}, "", 400, ""},
		{"missing operand", views.CalculationRequest{Operation: "sqrt"}, "", 400, ""},
	}
	for _, item := range tests {
		t.Run(item.name, func(t *testing.T) {
			body, err := json.Marshal(item.input)
			if err != nil {
				t.Fatal(err)
			}
			request := httptest.NewRequest("POST", "/api/calculate", strings.NewReader(string(body)))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			newHandler().ServeHTTP(response, request)
			if response.Code != item.status {
				t.Fatalf("status %d: %s", response.Code, response.Body.String())
			}
			if response.Header().Get("Content-Type") != "application/json" {
				t.Fatal("not JSON")
			}
			var data map[string]string
			if err := json.Unmarshal(response.Body.Bytes(), &data); err != nil {
				t.Fatal(err)
			}
			if item.status == 200 && data["result"] != item.result {
				t.Fatalf("result %q, want %q", data["result"], item.result)
			}
			if item.status != 200 && data["error"] == "" {
				t.Fatal("missing error message")
			}
			if item.message != "" && data["error"] != item.message {
				t.Fatalf("error %q", data["error"])
			}
		})
	}
}

func TestRequestValidation(t *testing.T) {
	for _, item := range []struct {
		name, body, contentType string
		status                  int
	}{
		{"malformed", `{"operation":`, "application/json", 400},
		{"unknown field", `{"operation":"sqrt","left":"9","extra":true}`, "application/json", 400},
		{"multiple objects", `{"operation":"sqrt","left":"9"} {}`, "application/json", 400},
		{"numeric operand", `{"operation":"sqrt","left":9}`, "application/json", 400},
		{"array", `[]`, "application/json", 400},
		{"null", `null`, "application/json", 400},
		{"content type", `{"operation":"sqrt","left":"9"}`, "text/plain", 415},
		{"too large", strings.Repeat(" ", 5000), "application/json", 413},
		{"large trailing data", `{"operation":"sqrt","left":"9"}` + strings.Repeat(" ", 5000), "application/json", 413},
	} {
		t.Run(item.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/api/calculate", strings.NewReader(item.body))
			request.Header.Set("Content-Type", item.contentType)
			response := httptest.NewRecorder()
			newHandler().ServeHTTP(response, request)
			if response.Code != item.status {
				t.Fatalf("status %d: %s", response.Code, response.Body.String())
			}
		})
	}
	response := httptest.NewRecorder()
	newHandler().ServeHTTP(response, httptest.NewRequest("GET", "/api/calculate", nil))
	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET returned %d", response.Code)
	}
}
