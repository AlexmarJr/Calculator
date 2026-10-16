package routes_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"calculator/backend/api/controllers"
	"calculator/backend/api/repositories"
	"calculator/backend/api/routes"
	"calculator/backend/api/services"
)

func TestRouteContract(t *testing.T) {
	handler := routes.NewRouter(controllers.NewCalculatorController(
		services.NewCalculatorService(repositories.NewMemoryHistory),
	))
	for _, test := range []struct {
		method, path string
		status       int
		allow        string
	}{
		{http.MethodGet, "/api/health", http.StatusOK, ""},
		{http.MethodPost, "/api/health", http.StatusMethodNotAllowed, "GET, HEAD"},
		{http.MethodGet, "/api/calculate", http.StatusMethodNotAllowed, "POST"},
		{http.MethodGet, "/api/calculator/action", http.StatusMethodNotAllowed, "POST"},
		{http.MethodGet, "/api/unknown", http.StatusNotFound, ""},
	} {
		t.Run(test.method+" "+test.path, func(t *testing.T) {
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, httptest.NewRequest(test.method, test.path, nil))
			if response.Code != test.status {
				t.Fatalf("status = %d, want %d", response.Code, test.status)
			}
			if got := response.Header().Get("Allow"); got != test.allow {
				t.Fatalf("Allow = %q, want %q", got, test.allow)
			}
			if test.status == http.StatusOK {
				var body struct {
					Status string `json:"status"`
				}
				if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
					t.Fatal(err)
				}
				if body.Status != "ok" || response.Header().Get("Content-Type") != "application/json" {
					t.Fatalf("invalid health response: %s", response.Body.String())
				}
			}
		})
	}
}
