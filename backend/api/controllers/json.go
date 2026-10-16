package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"

	"calculator/backend/api/models"
	"calculator/backend/api/views"
)

func writeServiceError(w http.ResponseWriter, err error) {
	var failure *models.Error
	if errors.As(err, &failure) {
		switch failure.Kind {
		case models.InvalidInput:
			views.Error(w, http.StatusBadRequest, failure.Message)
			return
		case models.UndefinedResult:
			views.Error(w, http.StatusUnprocessableEntity, failure.Message)
			return
		}
	}
	views.Error(w, http.StatusInternalServerError, "Internal server error.")
}

func decodeJSON(w http.ResponseWriter, r *http.Request, limit int64, target any, message string) bool {
	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || contentType != "application/json" {
		views.Error(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json.")
		return false
	}
	r.Body = http.MaxBytesReader(w, r.Body, limit)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeDecodeError(w, err, message)
		return false
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		writeDecodeError(w, err, "Send exactly one JSON object.")
		return false
	}
	return true
}

func writeDecodeError(w http.ResponseWriter, err error, message string) {
	var sizeError *http.MaxBytesError
	if errors.As(err, &sizeError) {
		views.Error(w, http.StatusRequestEntityTooLarge, "Request body is too large.")
		return
	}
	views.Error(w, http.StatusBadRequest, message)
}
