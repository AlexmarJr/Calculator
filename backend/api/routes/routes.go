package routes

import (
	"net/http"

	"calculator/backend/api/controllers"
)

// NewRouter registers the API endpoints and their allowed HTTP methods.
func NewRouter(calculator *controllers.CalculatorController) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", controllers.Health)
	mux.HandleFunc("POST /api/calculate", calculator.Calculate)
	mux.HandleFunc("POST /api/calculator/action", calculator.Action)
	return mux
}
