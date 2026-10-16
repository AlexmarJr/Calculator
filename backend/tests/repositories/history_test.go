package repositories_test

import (
	"testing"

	"calculator/backend/api/models"
	"calculator/backend/api/repositories"
)

func TestMemoryHistoryIsolationAndCopies(t *testing.T) {
	seed := []models.HistoryEntry{{ID: 1, Expression: "1 + 1", Result: "2"}}
	first := repositories.NewMemoryHistory(seed)
	second := repositories.NewMemoryHistory(seed)
	first.Save(models.HistoryEntry{
		ID:         2,
		Expression: "2 + 2",
		Result:     "4",
	})
	if len(second.List()) != 1 || len(seed) != 1 {
		t.Fatal("repository instances shared request state")
	}
	list := first.List()
	list[0].Result = "changed"
	found, ok := first.Find(2)
	if !ok || found.Result != "4" {
		t.Fatal("repository exposed mutable storage")
	}
	first.Save(models.HistoryEntry{
		ID:         2,
		Expression: "updated",
		Result:     "5",
	})
	if len(first.List()) != 2 || first.List()[0].Result != "5" {
		t.Fatal("save did not update an existing entry")
	}
}

func TestMemoryHistoryLimitAndClear(t *testing.T) {
	repository := repositories.NewMemoryHistory(nil)
	for id := int64(1); id <= 30; id++ {
		repository.Save(models.HistoryEntry{
			ID:         id,
			Expression: "1 + 1",
			Result:     "2",
		})
	}
	entries := repository.List()
	if len(entries) != models.HistoryLimit || entries[0].ID != 30 || entries[len(entries)-1].ID != 11 {
		t.Fatal("repository did not retain the 20 most recent entries")
	}
	if _, found := repository.Find(1); found {
		t.Fatal("expired entry was retained")
	}
	repository.Clear()
	if repository.List() == nil || len(repository.List()) != 0 {
		t.Fatal("history was not cleared")
	}
}
