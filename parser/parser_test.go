package parser

import (
	"testing"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/scanner"
	"github.com/shopspring/decimal"
)

func TestParseBinaryOperator(t *testing.T) {
	tests := []struct {
		input    string
		left     int
		operator string
		right    int
	}{
		{"1 + 5", 1, "+", 5},
		{"1 - 5", 1, "-", 5},
		{"1 * 5", 1, "*", 5},
		{"1 / 5", 1, "/", 5},
		{"1 == 5", 1, "==", 5},
		{"1 != 5", 1, "!=", 5},
		{"1 >= 5", 1, ">=", 5},
		{"1 <= 5", 1, "<=", 5},
		{"1 > 5", 1, ">", 5},
		{"1 < 5", 1, "<", 5},
	}

	for i, test := range tests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		expression := parser.Parse()

		binary, ok := expression.(*ast.Binary)

		if !ok {
			t.Fatalf("test[%v] result is not *ast.Binary, got=%T", i, expression)
		}

		if binary.Operator.Lexeme != test.operator {
			t.Errorf("binary operator value not %v, got=%v", test.operator, binary.Operator.Lexeme)
		}

		verifyNumberLiteral(binary.Left, test.left, t)
		verifyNumberLiteral(binary.Right, test.right, t)
	}
}

func TestParseBooleans(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, test := range tests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		expression := parser.Parse()

		literal, ok := expression.(*ast.Literal)

		if !ok {
			t.Fatalf("result is not *ast.Literal, got=%T", expression)
		}

		value, ok := literal.Value.(bool)

		if !ok {
			t.Fatalf("Literal.Value type not bool, got=%T", value)
		}

		if value != test.expected {
			t.Errorf("literal value not %v, got=%v", test.expected, value)
		}
	}
}

func TestParseGroupedExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"(5)", 5},
		{"(3.14)", 3.14},
	}

	for _, test := range tests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		expression := parser.Parse()

		grouping, ok := expression.(*ast.Grouping)

		if !ok {
			t.Fatalf("expression is not *ast.Grouping, got=%T", expression)
		}

		verifyNumberLiteral(grouping.Expression, test.expected, t)
	}
}

func TestParseNull(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"null"},
	}

	for _, test := range tests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		expression := parser.Parse()

		literal, ok := expression.(*ast.Literal)

		if !ok {
			t.Fatalf("result is not *ast.Literal, got=%T", expression)
		}

		if literal.Value != nil {
			t.Errorf("literal value not %v, got=%v", nil, literal.Value)
		}
	}
}

func TestParseNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"5", 5},
		{"3.14", 3.14},
	}

	for _, test := range tests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		expression := parser.Parse()

		verifyNumberLiteral(expression, test.expected, t)
	}
}

func TestParseStrings(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"\"hello\"", "hello"},
		{"\"world\"", "world"},
	}

	for _, test := range tests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		expression := parser.Parse()

		literal, ok := expression.(*ast.Literal)

		if !ok {
			t.Fatalf("result is not *ast.Literal, got=%T", expression)
		}

		value, ok := literal.Value.(string)

		if !ok {
			t.Fatalf("Literal.Value type not string, got=%T", value)
		}

		if value != test.expected {
			t.Errorf("literal value not %v, got=%v", test.expected, value)
		}
	}
}

func TestParseUnaryOperators(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		right    interface{}
	}{
		{"!true", "!", true},
		{"!false", "!", false},
		{"-3.14", "-", 3.14},
	}

	for _, test := range tests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		expression := parser.Parse()

		unary, ok := expression.(*ast.Unary)

		if !ok {
			t.Fatalf("result is not *ast.Unary, got=%T", expression)
		}

		if unary.Operator.Lexeme != test.operator {
			t.Errorf("binary operator value not %v, got=%v", test.operator, unary.Operator.Lexeme)
		}

		right, ok := unary.Right.(*ast.Literal)

		if !ok {
			t.Fatalf("unary right is not *ast.Literal, got=%T", right)
		}

		value, ok := right.Value.(decimal.Decimal)

		if !ok {
			value, ok := right.Value.(bool)

			if !ok {
				t.Fatalf("unary right type is not decimal.Decimal or bool, got=%T", value)
			} else if value != test.right.(bool) {
				t.Errorf("unary right value not %v, got=%v", test.right, value)
			}
		} else {
			number := decimal.NewFromFloat(test.right.(float64))

			if !number.Equal(value) {
				t.Errorf("unary right value not %v, got=%v", number, value)
			}
		}
	}
}

// =============================================================================
// Helper methods

func verifyNumberLiteral(expression ast.ExpressionNode, expected interface{}, t *testing.T) {
	literal, ok := expression.(*ast.Literal)

	if !ok {
		t.Fatalf("result is not *ast.Literal, got=%T", expression)
	}

	check, ok := expected.(int)

	if ok {
		expected = decimal.NewFromInt(int64(check))
	} else {
		check, ok := expected.(float64)

		if ok {
			expected = decimal.NewFromFloat(check)
		} else {
			t.Fatalf("Expected either an int64 or float64, got=%T", expected)
		}
	}

	equals := expected.(decimal.Decimal).Equal(literal.Value.(decimal.Decimal))

	if !equals {
		t.Errorf("literal value not %v, got=%v", expected, literal.Value)
	}
}
