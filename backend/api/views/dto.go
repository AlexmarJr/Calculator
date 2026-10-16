package views

import (
	"encoding/json"
	"errors"
	"strings"

	"calculator/backend/api/models"
)

type CalculationRequest struct {
	Operation string  `json:"operation"`
	Left      string  `json:"left"`
	Right     *string `json:"right,omitempty"`
	Base      *string `json:"base,omitempty"`
}

type HistoryEntry struct {
	ID         int64  `json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

type CalculatorState struct {
	Display         string         `json:"display"`
	Stored          *string        `json:"stored"`
	Operator        *string        `json:"operator"`
	AwaitingOperand bool           `json:"awaitingOperand"`
	JustCalculated  bool           `json:"justCalculated"`
	Expression      string         `json:"expression"`
	OperandLabel    *string        `json:"operandLabel"`
	Error           *string        `json:"error"`
	History         []HistoryEntry `json:"history"`
	NextID          int64          `json:"nextId"`
}

type CalculatorAction struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value,omitempty"`
}

type ActionRequest struct {
	State  *CalculatorState  `json:"state"`
	Action *CalculatorAction `json:"action"`
}

func (input CalculationRequest) ToModel() models.Calculation {
	return models.Calculation{
		Operation: input.Operation,
		Left:      input.Left,
		Right:     input.Right,
		Base:      input.Base,
	}
}

func (state CalculatorState) ToModel() models.State {
	history := make([]models.HistoryEntry, len(state.History))
	for index, entry := range state.History {
		history[index] = models.HistoryEntry{
			ID:         entry.ID,
			Expression: entry.Expression,
			Result:     entry.Result,
		}
	}
	return models.State{
		Display:         state.Display,
		Stored:          state.Stored,
		Operator:        state.Operator,
		AwaitingOperand: state.AwaitingOperand,
		JustCalculated:  state.JustCalculated,
		Expression:      state.Expression,
		OperandLabel:    state.OperandLabel,
		Error:           state.Error,
		History:         history,
		NextID:          state.NextID,
	}
}

func StateFromModel(state models.State) CalculatorState {
	history := make([]HistoryEntry, len(state.History))
	for index, entry := range state.History {
		history[index] = HistoryEntry{
			ID:         entry.ID,
			Expression: entry.Expression,
			Result:     entry.Result,
		}
	}
	return CalculatorState{
		Display:         state.Display,
		Stored:          state.Stored,
		Operator:        state.Operator,
		AwaitingOperand: state.AwaitingOperand,
		JustCalculated:  state.JustCalculated,
		Expression:      state.Expression,
		OperandLabel:    state.OperandLabel,
		Error:           state.Error,
		History:         history,
		NextID:          state.NextID,
	}
}

func (action CalculatorAction) ToModel() (models.Action, error) {
	command := models.Action{Type: action.Type, HasValue: len(action.Value) != 0}
	switch action.Type {
	case "digit", "operator":
		if err := json.Unmarshal(action.Value, &command.Value); err != nil {
			return command, errors.New("Action value must be a string.")
		}
	case "restore":
		var entry HistoryEntry
		decoder := json.NewDecoder(strings.NewReader(string(action.Value)))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&entry); err != nil {
			return command, errors.New("Invalid history entry.")
		}
		command.HistoryID = entry.ID
	}
	return command, nil
}
