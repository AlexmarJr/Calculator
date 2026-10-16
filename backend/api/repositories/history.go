package repositories

import "calculator/backend/api/models"

// Each request receives its own repository seeded from the submitted history.
type HistoryRepository interface {
	List() []models.HistoryEntry
	Find(id int64) (models.HistoryEntry, bool)
	Save(entry models.HistoryEntry)
	Clear()
}

type HistoryFactory func(initial []models.HistoryEntry) HistoryRepository

type memoryHistory struct {
	entries []models.HistoryEntry
}

func NewMemoryHistory(initial []models.HistoryEntry) HistoryRepository {
	repository := &memoryHistory{entries: append([]models.HistoryEntry{}, initial...)}
	if len(repository.entries) > models.HistoryLimit {
		repository.entries = repository.entries[:models.HistoryLimit]
	}
	return repository
}

func (r *memoryHistory) List() []models.HistoryEntry {
	return append([]models.HistoryEntry{}, r.entries...)
}

func (r *memoryHistory) Find(id int64) (models.HistoryEntry, bool) {
	for _, entry := range r.entries {
		if entry.ID == id {
			return entry, true
		}
	}
	return models.HistoryEntry{}, false
}

func (r *memoryHistory) Save(entry models.HistoryEntry) {
	for index, current := range r.entries {
		if current.ID == entry.ID {
			r.entries = append(r.entries[:index], r.entries[index+1:]...)
			break
		}
	}
	r.entries = append([]models.HistoryEntry{entry}, r.entries...)
	if len(r.entries) > models.HistoryLimit {
		r.entries = r.entries[:models.HistoryLimit]
	}
}

func (r *memoryHistory) Clear() {
	r.entries = []models.HistoryEntry{}
}
