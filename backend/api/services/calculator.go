package services

import (
	"calculator/backend/api/models"
	"calculator/backend/api/repositories"
)

type CalculatorService struct {
	historyFactory repositories.HistoryFactory
}

func NewCalculatorService(historyFactory repositories.HistoryFactory) *CalculatorService {
	return &CalculatorService{historyFactory: historyFactory}
}

func (s *CalculatorService) Calculate(input models.Calculation) (string, error) {
	return Calculate(input)
}

func (s *CalculatorService) Apply(state models.State, action models.Action) (models.State, error) {
	if err := ValidateState(state); err != nil {
		return state, err
	}
	history := s.historyFactory(state.History)
	state.History = history.List()
	next, err := applyAction(state, action, history)
	if err != nil {
		return state, err
	}
	next.History = history.List()
	if err := ValidateState(next); err != nil {
		return state, err
	}
	return next, nil
}
