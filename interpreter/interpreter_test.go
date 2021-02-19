package interpreter

import (
	"testing"

	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
	"github.com/shopspring/decimal"
)

func TestEvaluateLiteral(t *testing.T) {
	tests := []struct {
		literal  string
		expected interface{}
	}{
		{"5", 5},
	}

	for _, test := range tests {
		scanner := scanner.New(test.literal)
		tokens := scanner.ScanTokens()
		parser := parser.New(tokens)
		expression := parser.Parse()

		result := Evaluate(expression)

		verifyLiteralValue(result, test.expected, t)
	}
}

// =============================================================================
// Helper methods

func verifyLiteralValue(literal interface{}, expected interface{}, t *testing.T) {
	switch result := literal.(type) {
	case decimal.Decimal:
		verifyNumberValue(result, expected, t)
	// case bool:
	// 	verifyBooleanValue(result, expected, t)
	// case string:
	// 	verifyStringValue(result, expected, t)
	default:
		t.Fatalf("Unsupported literal type, expected float64, bool, or string, got=%T", result)
	}
}

func verifyNumberValue(value decimal.Decimal, expected interface{}, t *testing.T) {
	check, ok := expected.(int)

	if ok {
		expected = decimal.NewFromInt(int64(check))
	} else {
		check, ok := expected.(float64)

		if ok {
			expected = decimal.NewFromFloat(check)
		} else {
			t.Fatalf("Expected either an int or float64, got=%T", expected)
		}
	}

	equals := expected.(decimal.Decimal).Equal(value)

	if !equals {
		t.Errorf("Numbers are not equal, expected %v, got=%v", expected, value)
	}
}
