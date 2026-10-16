package services

import (
	"math"
	"regexp"
	"strconv"

	"calculator/backend/api/models"
)

var numberPattern = regexp.MustCompile(`^[+-]?(?:\d+(?:\.\d*)?|\.\d+)(?:[eE][+-]?\d+)?$`)

func parseOperand(value string) (float64, error) {
	if len(value) > 64 || !numberPattern.MatchString(value) {
		return 0, models.Invalid("Operands must be decimal strings with a dot as the decimal separator.")
	}
	number, err := strconv.ParseFloat(value, 64)
	if err != nil || math.IsNaN(number) || math.IsInf(number, 0) {
		return 0, models.Invalid("Operands must be finite numbers.")
	}
	return number, nil
}

func Calculate(input models.Calculation) (string, error) {
	binary := false
	switch input.Operation {
	case "add", "subtract", "multiply", "divide", "power":
		binary = true
	case "sqrt", "percent":
	default:
		return "", models.Invalid("Unsupported operation.")
	}
	if binary && input.Right == nil {
		return "", models.Invalid("This operation requires a right operand.")
	}
	if !binary && input.Right != nil {
		return "", models.Invalid("This operation takes only a left operand.")
	}
	if input.Base != nil && input.Operation != "percent" {
		return "", models.Invalid("A base is only supported for percentages.")
	}
	left, err := parseOperand(input.Left)
	if err != nil {
		return "", err
	}
	var right, base float64
	if input.Right != nil {
		right, err = parseOperand(*input.Right)
		if err != nil {
			return "", err
		}
	}
	if input.Base != nil {
		base, err = parseOperand(*input.Base)
		if err != nil {
			return "", err
		}
	}
	var result float64
	switch input.Operation {
	case "add":
		result = left + right
	case "subtract":
		result = left - right
	case "multiply":
		result = left * right
	case "divide":
		if right == 0 {
			return "", models.Undefined("Cannot divide by zero.")
		}
		result = left / right
	case "power":
		result = math.Pow(left, right)
	case "sqrt":
		if left < 0 {
			return "", models.Undefined("Square root needs a non-negative number.")
		}
		result = math.Sqrt(left)
	case "percent":
		result = left / 100
		if input.Base != nil {
			result *= base
		}
	}
	if math.IsNaN(result) {
		return "", models.Undefined("The result is not a real number.")
	}
	if math.IsInf(result, 0) {
		return "", models.Undefined("The result is too large.")
	}
	// Round results to 15 significant digits, then use ordinary decimal notation
	// for everyday magnitudes. This also removes floating-point noise.
	rounded, err := strconv.ParseFloat(strconv.FormatFloat(result, 'g', 15, 64), 64)
	if err != nil {
		return "", models.Undefined("The result is too large.")
	}
	if rounded == 0 {
		return "0", nil
	}
	if math.Abs(rounded) >= 1e-6 && math.Abs(rounded) < 1e21 {
		return strconv.FormatFloat(rounded, 'f', -1, 64), nil
	}
	return strconv.FormatFloat(rounded, 'g', -1, 64), nil
}
