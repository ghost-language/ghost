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
		expression, _ := parser.expression()

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
		expression, _ := parser.expression()

		literal, ok := expression.(*ast.Boolean)

		if !ok {
			t.Fatalf("result is not *ast.Boolean, got=%T", expression)
		}

		value := literal.Value

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
		expression, _ := parser.expression()

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
		expression, _ := parser.expression()

		_, ok := expression.(*ast.Null)

		if !ok {
			t.Fatalf("result is not *ast.Null, got=%T", expression)
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
		expression, _ := parser.expression()

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
		expression, _ := parser.expression()

		literal, ok := expression.(*ast.String)

		if !ok {
			t.Fatalf("result is not *ast.String, got=%T", expression)
		}

		value := literal.Value

		if value != test.expected {
			t.Errorf("string value not %v, got=%v", test.expected, value)
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
	}

	for _, test := range tests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		expression, _ := parser.expression()

		unary, ok := expression.(*ast.Unary)

		if !ok {
			t.Fatalf("result is not *ast.Unary, got=%T", expression)
		}

		if unary.Operator.Lexeme != test.operator {
			t.Errorf("binary operator value not %v, got=%v", test.operator, unary.Operator.Lexeme)
		}

		right, ok := unary.Right.(*ast.Boolean)

		if !ok {
			t.Fatalf("unary right is not *ast.Boolean, got=%T", right)
		}

		if right.Value != test.right.(bool) {
			t.Errorf("unary right value not %v, got=%v", test.right, right.Value)
		}
	}
}

func TestParseExpressionStatement(t *testing.T) {
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
		statements := parser.Parse()

		if len(statements) != 1 {
			t.Fatalf("Expected 1 statement, got=%v", len(statements))
		}

		expression, ok := statements[0].(*ast.Expression)

		if !ok {
			t.Fatalf("Expected *ast.Expression, got=%T", statements[0])
		}

		verifyNumberLiteral(expression.Expression, test.expected, t)
	}
}

func TestParseVariables(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"foo"},
		{"bar"},
	}

	for _, test := range tests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		statements := parser.Parse()

		_, ok := statements[0].(*ast.Expression)

		if !ok {
			t.Fatalf("Expected *ast.Expression, got=%T", statements[0])
		}
	}
}

// =============================================================================
// Helper methods

func verifyNumberLiteral(expression ast.ExpressionNode, expected interface{}, t *testing.T) {
	number, ok := expression.(*ast.Number)

	if !ok {
		t.Fatalf("result is not *ast.Number, got=%T", expression)
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

	equals := expected.(decimal.Decimal).Equal(number.Value)

	if !equals {
		t.Errorf("number value not %v, got=%v", expected, number.Value)
	}
}
