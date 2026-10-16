package controllers

import (
	"net/http"

	"calculator/backend/api/views"
)

func Health(w http.ResponseWriter, r *http.Request) {
	views.Health(w)
}
