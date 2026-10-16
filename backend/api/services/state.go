package services

import (
	"strings"

	"calculator/backend/api/models"
	"calculator/backend/api/repositories"
)

var operationNames = map[string]string{
	"+": "add",
	"-": "subtract",
	"*": "multiply",
	"/": "divide",
	"^": "power",
}

var operationSymbols = map[string]string{
	"+": "+",
	"-": "−",
	"*": "×",
	"/": "÷",
	"^": "^",
}

func textPointer(value string) *string {
	return &value
}

func resetCalculator(state models.State) models.State {
	fresh := models.InitialState()
	fresh.History, fresh.NextID = state.History, state.NextID
	return fresh
}

func calculationError(state models.State, err error) models.State {
	state = resetCalculator(state)
	state.JustCalculated = true
	state.Error = textPointer(err.Error())
	return state
}

func beginEntry(state models.State, value string) models.State {
	state.Display = value
	state.AwaitingOperand, state.JustCalculated = false, false
	state.Error, state.OperandLabel = nil, nil
	if state.Operator == nil {
		state.Expression = ""
	}
	return state
}

func recordResult(
	state models.State,
	history repositories.HistoryRepository,
	expression, result string,
) models.State {
	history.Save(models.HistoryEntry{
		ID:         state.NextID,
		Expression: expression,
		Result:     result,
	})
	state.NextID++
	state.History = history.List()
	return state
}

func ValidateState(state models.State) error {
	if _, err := parseOperand(state.Display); err != nil {
		return err
	}
	if state.NextID < 1 || state.NextID > 9007199254740970 || len(state.History) > models.HistoryLimit ||
		len(state.Expression) > 1024 || (state.OperandLabel != nil && len(*state.OperandLabel) > 1024) ||
		(state.Error != nil && len(*state.Error) > 512) {
		return models.Invalid("Invalid calculator state.")
	}
	if (state.Stored == nil) != (state.Operator == nil) || (state.AwaitingOperand && state.Operator == nil) {
		return models.Invalid("Invalid pending operation.")
	}
	if state.Operator != nil {
		if _, ok := operationNames[*state.Operator]; !ok {
			return models.Invalid("Unsupported operator.")
		}
		if _, err := parseOperand(*state.Stored); err != nil {
			return err
		}
	}
	ids := make(map[int64]bool)
	for _, item := range state.History {
		if item.ID < 1 || item.ID >= state.NextID || ids[item.ID] || len(item.Expression) > 1024 {
			return models.Invalid("Invalid calculator history.")
		}
		ids[item.ID] = true
		if _, err := parseOperand(item.Result); err != nil {
			return err
		}
	}
	return nil
}

func applyAction(
	state models.State,
	action models.Action,
	history repositories.HistoryRepository,
) (models.State, error) {
	value := action.Value
	switch action.Type {
	case "digit":
		if !action.HasValue {
			return state, models.Invalid("models.Action value must be a string.")
		}
		if len(value) != 1 || value[0] < '0' || value[0] > '9' {
			return state, models.Invalid("Digit must be between 0 and 9.")
		}
	case "operator":
		if !action.HasValue {
			return state, models.Invalid("models.Action value must be a string.")
		}
		if _, ok := operationNames[value]; !ok {
			return state, models.Invalid("Unsupported operator.")
		}
	case "restore":
		if item, found := history.Find(action.HistoryID); found {
			state = resetCalculator(state)
			state.Display, state.Expression, state.JustCalculated = item.Result, item.Expression+" =", true
			return state, nil
		}
		return state, models.Invalid("History entry was not found.")
	case "decimal", "equals", "clear", "backspace", "sign", "percent", "sqrt", "clear-history":
		if action.HasValue {
			return state, models.Invalid("This action does not take a value.")
		}
	default:
		return state, models.Invalid("Unsupported calculator action.")
	}
	switch action.Type {
	case "clear":
		return resetCalculator(state), nil
	case "clear-history":
		history.Clear()
		state.History = history.List()
	case "digit":
		if state.Error != nil {
			return beginEntry(resetCalculator(state), value), nil
		}
		if state.AwaitingOperand || state.JustCalculated || state.Error != nil {
			return beginEntry(state, value), nil
		}
		digits := 0
		for _, character := range state.Display {
			if character >= '0' && character <= '9' {
				digits++
			}
		}
		if digits >= 15 {
			return state, nil
		}
		switch state.Display {
		case "0":
			state.Display = value
		case "-0":
			state.Display = "-" + value
		default:
			state.Display += value
		}
		state.OperandLabel = nil
	case "decimal":
		if state.Error != nil {
			return beginEntry(resetCalculator(state), "0."), nil
		}
		if state.AwaitingOperand || state.JustCalculated || state.Error != nil {
			return beginEntry(state, "0."), nil
		}
		if !strings.Contains(state.Display, ".") {
			state.Display += "."
			state.OperandLabel = nil
		}
	case "backspace":
		if state.Error != nil {
			return resetCalculator(state), nil
		}
		if state.AwaitingOperand || state.JustCalculated {
			return state, nil
		}
		if len(state.Display) <= 1 || (len(state.Display) == 2 && state.Display[0] == '-') {
			state.Display = "0"
		} else {
			state.Display = state.Display[:len(state.Display)-1]
		}
		state.OperandLabel = nil
	case "sign":
		if state.Error != nil {
			return state, nil
		}
		if state.AwaitingOperand {
			return beginEntry(state, "-0"), nil
		}
		if strings.HasPrefix(state.Display, "-") {
			state.Display = state.Display[1:]
		} else {
			state.Display = "-" + state.Display
		}
		state.OperandLabel = nil
	case "operator":
		if state.Error != nil {
			return state, nil
		}
		if value == "^" && state.Operator != nil && *state.Operator == "^" {
			state.JustCalculated = state.AwaitingOperand || state.JustCalculated
			state.Stored, state.Operator, state.OperandLabel = nil, nil, nil
			state.Expression, state.AwaitingOperand = "", false
			return state, nil
		}
		display := state.Display
		if state.Operator != nil {
			if state.AwaitingOperand {
				display = *state.Stored
			} else {
				result, err := Calculate(models.Calculation{
					Operation: operationNames[*state.Operator],
					Left:      *state.Stored,
					Right:     &state.Display,
				})
				if err != nil {
					return calculationError(state, err), nil
				}
				display = result
			}
		}
		state.Display, state.Stored, state.Operator = display, textPointer(display), textPointer(value)
		state.AwaitingOperand, state.JustCalculated = true, false
		state.OperandLabel, state.Expression = nil, display+" "+operationSymbols[value]
	case "equals":
		if state.Error != nil || state.Operator == nil || state.AwaitingOperand {
			return state, nil
		}
		label := state.Display
		if state.OperandLabel != nil {
			label = *state.OperandLabel
		}
		expression := *state.Stored + " " + operationSymbols[*state.Operator] + " " + label
		result, err := Calculate(models.Calculation{
			Operation: operationNames[*state.Operator],
			Left:      *state.Stored,
			Right:     &state.Display,
		})
		if err != nil {
			return calculationError(state, err), nil
		}
		state.Display, state.Expression = result, expression+" ="
		state.Stored, state.Operator, state.OperandLabel = nil, nil, nil
		state.JustCalculated, state.AwaitingOperand = true, false
		state = recordResult(state, history, expression, result)
	case "sqrt", "percent":
		if state.Error != nil || state.AwaitingOperand {
			return state, nil
		}
		label := state.Display
		if state.OperandLabel != nil {
			label = *state.OperandLabel
		}
		input := models.Calculation{
			Operation: action.Type,
			Left:      state.Display,
		}
		relative := action.Type == "percent" &&
			state.Operator != nil &&
			(*state.Operator == "+" || *state.Operator == "-")
		if relative {
			input.Base = state.Stored
		}
		result, err := Calculate(input)
		if err != nil {
			return calculationError(state, err), nil
		}
		operandLabel := label + "%"
		if action.Type == "sqrt" {
			operandLabel = "√(" + label + ")"
		}
		expression := operandLabel
		if relative {
			expression += " of " + *state.Stored
		}
		state.Display, state.OperandLabel, state.JustCalculated = result, textPointer(operandLabel), true
		if state.Operator == nil {
			state.Expression = expression + " ="
		}
		state = recordResult(state, history, expression, result)
	}
	return state, nil
}
