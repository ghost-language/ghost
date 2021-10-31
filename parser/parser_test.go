package parser

import (
	"testing"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/scanner"
)

func TestBoolean(t *testing.T) {
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

func TestIdentifier(t *testing.T) {
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

func TestNumber(t *testing.T) {
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

func TestString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello world"`, `"hello world"`},
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
