package services_test

import (
	"errors"
	"testing"

	"calculator/backend/api/models"
	"calculator/backend/api/services"
)

func value(text string) *string {
	return &text
}

func TestCalculate(t *testing.T) {
	for _, item := range []struct {
		input    models.Calculation
		expected string
	}{
		{
			models.Calculation{
				Operation: "add",
				Left:      "8",
				Right:     value("3"),
			},
			"11",
		},
		{
			models.Calculation{
				Operation: "subtract",
				Left:      "8",
				Right:     value("3"),
			},
			"5",
		},
		{
			models.Calculation{
				Operation: "multiply",
				Left:      "8",
				Right:     value("3"),
			},
			"24",
		},
		{
			models.Calculation{
				Operation: "divide",
				Left:      "8",
				Right:     value("3"),
			},
			"2.66666666666667",
		},
		{
			models.Calculation{
				Operation: "divide",
				Left:      "1",
				Right:     value("3"),
			},
			"0.333333333333333",
		},
		{
			models.Calculation{
				Operation: "add",
				Left:      "123456789012345",
				Right:     value("0"),
			},
			"123456789012345",
		},
		{
			models.Calculation{
				Operation: "power",
				Left:      "2",
				Right:     value("10"),
			},
			"1024",
		},
		{
			models.Calculation{
				Operation: "sqrt",
				Left:      "81",
			},
			"9",
		},
		{
			models.Calculation{
				Operation: "percent",
				Left:      "10",
				Base:      value("200"),
			},
			"20",
		},
		{
			models.Calculation{
				Operation: "add",
				Left:      "0.1",
				Right:     value("0.2"),
			},
			"0.3",
		},
	} {
		result, err := services.Calculate(item.input)
		if err != nil || result != item.expected {
			t.Fatalf("%+v: result %q, error %v", item.input, result, err)
		}
	}
}

func TestBusinessErrorKinds(t *testing.T) {
	for _, item := range []struct {
		input models.Calculation
		kind  models.ErrorKind
	}{
		{
			models.Calculation{
				Operation: "unknown",
				Left:      "9",
			},
			models.InvalidInput,
		},
		{
			models.Calculation{
				Operation: "sqrt",
				Left:      "1.234,56",
			},
			models.InvalidInput,
		},
		{
			models.Calculation{
				Operation: "add",
				Left:      "9",
			},
			models.InvalidInput,
		},
		{
			models.Calculation{
				Operation: "divide",
				Left:      "9",
				Right:     value("0"),
			},
			models.UndefinedResult,
		},
		{
			models.Calculation{
				Operation: "sqrt",
				Left:      "-9",
			},
			models.UndefinedResult,
		},
		{
			models.Calculation{
				Operation: "power",
				Left:      "-2",
				Right:     value("0.5"),
			},
			models.UndefinedResult,
		},
		{
			models.Calculation{
				Operation: "power",
				Left:      "99",
				Right:     value("999"),
			},
			models.UndefinedResult,
		},
	} {
		_, err := services.Calculate(item.input)
		var failure *models.Error
		if !errors.As(err, &failure) || failure.Kind != item.kind {
			t.Fatalf("unexpected business error: %v", err)
		}
	}
}
