package services_test

import (
	"testing"

	"calculator/backend/api/models"
	"calculator/backend/api/repositories"
	"calculator/backend/api/services"
)

type recordingHistory struct {
	repositories.HistoryRepository
	saves, finds, clears int
}

func (r *recordingHistory) Save(entry models.HistoryEntry) {
	r.saves++
	r.HistoryRepository.Save(entry)
}

func (r *recordingHistory) Find(id int64) (models.HistoryEntry, bool) {
	r.finds++
	return r.HistoryRepository.Find(id)
}

func (r *recordingHistory) Clear() {
	r.clears++
	r.HistoryRepository.Clear()
}

func TestServiceUsesHistoryRepository(t *testing.T) {
	var current *recordingHistory
	service := services.NewCalculatorService(func(initial []models.HistoryEntry) repositories.HistoryRepository {
		current = &recordingHistory{HistoryRepository: repositories.NewMemoryHistory(initial)}
		return current
	})
	state := models.InitialState()
	state.Display = "9"
	state, err := service.Apply(state, models.Action{Type: "sqrt"})
	if err != nil || current.saves != 1 || len(state.History) != 1 {
		t.Fatalf("result was not stored: %v", err)
	}
	state, err = service.Apply(state, models.Action{Type: "restore", HistoryID: state.History[0].ID})
	if err != nil || current.finds != 1 || state.Display != "3" {
		t.Fatalf("history was not retrieved: %v", err)
	}
	state, err = service.Apply(state, models.Action{Type: "clear-history"})
	if err != nil || current.clears != 1 || len(state.History) != 0 {
		t.Fatalf("history was not cleared: %v", err)
	}
}
