package parser

import (
	"testing"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/scanner"
)

func TestAssignStatement(t *testing.T) {
	tests := []struct {
		input      string
		identifier string
		value      int64
	}{
		{`a := 5`, "a", 5},
	}

	for _, tt := range tests {
		scanner := scanner.New(tt.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		statements := parser.Parse()

		if len(statements) != 1 {
			t.Fatalf("statements does not contain 1 statement. got=%d", len(statements))
		}

		assign, ok := statements[0].(*ast.Assign)

		if !ok {
			t.Fatalf("statements[0] is not ast.Assign. got=%T", statements[0])
		}

		if assign.Token.Lexeme != tt.identifier {
			t.Fatalf("assign.Token is not '%s'. got=%s", tt.identifier, assign.Token.Lexeme)
		}

		if !isNumberLiteral(t, assign.Value, tt.value) {
			return
		}
	}
}

func TestBooleanLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		scanner := scanner.New(tt.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		statements := parser.Parse()

		if len(statements) != 1 {
			t.Fatalf("statements does not contain 1 statement. got=%d", len(statements))
		}

		statement, ok := statements[0].(*ast.Expression)

		if !ok {
			t.Fatalf("statements[0] is not ast.Expression. got=%T", statements[0])
		}

		boolean, ok := statement.Expression.(*ast.Boolean)

		if !ok {
			t.Fatalf("statement is not ast.Boolean. got=%T", statement.Expression)
		}

		if boolean.Value != tt.expected {
			t.Fatalf("boolean.Value is not '%t'. got=%t", tt.expected, boolean.Value)
		}
	}
}

func TestIdentifierLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"foobar", "foobar"},
	}

	for _, tt := range tests {
		scanner := scanner.New(tt.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		statements := parser.Parse()

		if len(statements) != 1 {
			t.Fatalf("statements does not contain 1 statement. got=%d", len(statements))
		}

		statement, ok := statements[0].(*ast.Expression)

		if !ok {
			t.Fatalf("statements[0] is not ast.Expression. got=%T", statements[0])
		}

		identifier, ok := statement.Expression.(*ast.Identifier)

		if !ok {
			t.Fatalf("statement is not ast.Identifier. got=%T", statement.Expression)
		}

		if identifier.Value != tt.expected {
			t.Fatalf("identifier.Value is not '%s'. got=%s", tt.expected, identifier.Value)
		}
	}
}

func TestNumberLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5", "5"},
		{"3.14", "3.14"},
		{"5e10", "50000000000"},
	}

	for _, tt := range tests {
		scanner := scanner.New(tt.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		statements := parser.Parse()

		if len(statements) != 1 {
			t.Fatalf("statements does not contain 1 statement. got=%d", len(statements))
		}

		statement, ok := statements[0].(*ast.Expression)

		if !ok {
			t.Fatalf("statements[0] is not ast.Expression. got=%T", statements[0])
		}

		number, ok := statement.Expression.(*ast.Number)

		if !ok {
			t.Fatalf("statement is not ast.Number. got=%T", statement.Expression)
		}

		if number.Value.String() != tt.expected {
			t.Fatalf("number.Value is not '%s'. got=%s", tt.expected, number.Value.String())
		}
	}
}

func TestStringLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello world"`, `hello world`},
	}

	for _, tt := range tests {
		scanner := scanner.New(tt.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		statements := parser.Parse()

		if len(statements) != 1 {
			t.Fatalf("statements does not contain 1 statement. got=%d", len(statements))
		}

		statement, ok := statements[0].(*ast.Expression)

		if !ok {
			t.Fatalf("statements[0] is not ast.Expression. got=%T", statements[0])
		}

		str, ok := statement.Expression.(*ast.String)

		if !ok {
			t.Fatalf("statement is not ast.String. got=%T", statement.Expression)
		}

		if str.Value != tt.expected {
			t.Fatalf("string.Value is not '%s'. got=%s", tt.expected, str.Value)
		}
	}
}

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		number   int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range tests {
		scanner := scanner.New(tt.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		statements := parser.Parse()

		if len(statements) != 1 {
			t.Fatalf("statements does not contain 1 statement. got=%d", len(statements))
		}

		statement, ok := statements[0].(*ast.Expression)

		if !ok {
			t.Fatalf("statements[0] is not ast.Expression. got=%T", statements[0])
		}

		prefix, ok := statement.Expression.(*ast.Prefix)

		if !ok {
			t.Fatalf("statement is not ast.Prefix. got=%T", statement.Expression)
		}

		if prefix.Operator != tt.operator {
			t.Fatalf("prefix.Operator is not '%s'. got=%s", tt.operator, prefix.Operator)
		}

		if !isNumberLiteral(t, prefix.Right, tt.number) {
			return
		}
	}
}

// =============================================================================
// Helper methods

func isNumberLiteral(t *testing.T, expression ast.ExpressionNode, value int64) bool {
	number, ok := expression.(*ast.Number)

	if !ok {
		t.Errorf("expression is not ast.Number. got=%T", expression)
	}

	if number.Value.IntPart() != value {
		t.Errorf("number.Value is not %d. got=%d", value, number.Value.IntPart())
	}

	return true
}
