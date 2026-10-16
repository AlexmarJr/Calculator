package views

import (
	"encoding/json"
	"net/http"

	"calculator/backend/api/models"
)

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func Error(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func Health(w http.ResponseWriter) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func Calculation(w http.ResponseWriter, result string) {
	writeJSON(w, http.StatusOK, map[string]string{"result": result})
}

func State(w http.ResponseWriter, state models.State) {
	writeJSON(w, http.StatusOK, map[string]CalculatorState{"state": StateFromModel(state)})
}
