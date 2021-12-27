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
		program := parser.Parse()

		failIfParserHasErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		assign, ok := program.Statements[0].(*ast.Assign)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.Assign. got=%T", program.Statements[0])
		}

		if assign.Name.Value != tt.identifier {
			t.Fatalf("assign.Name is not '%s'. got=%s", tt.identifier, assign.Name.Value)
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
		program := parser.Parse()

		failIfParserHasErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.Expression)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.Expression. got=%T", program.Statements[0])
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
		program := parser.Parse()

		failIfParserHasErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.Expression)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.Expression. got=%T", program.Statements[0])
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

func TestIfExpressions(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"if (true) { print(true) }"},
	}

	for _, tt := range tests {
		scanner := scanner.New(tt.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		program := parser.Parse()

		failIfParserHasErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.Expression)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.Expression. got=%T", program.Statements[0])
		}

		ifExpression, ok := statement.Expression.(*ast.If)

		if !ok {
			t.Fatalf("statement is not ast.If. got=%T", statement.Expression)
		}

		_ = ifExpression
	}
}

func TestInfixExpressions(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
	}

	for _, tt := range tests {
		scanner := scanner.New(tt.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		program := parser.Parse()

		failIfParserHasErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.Expression)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.Expression. got=%T", program.Statements[0])
		}

		infix, ok := statement.Expression.(*ast.Infix)

		if !ok {
			t.Fatalf("statement is not ast.Infix. got=%T", statement.Expression)
		}

		if infix.Operator != tt.operator {
			t.Fatalf("infix.Operator is not '%s'. got=%s", tt.operator, infix.Operator)
		}

		if !isNumberLiteral(t, infix.Right, tt.rightValue) {
			return
		}

		if !isNumberLiteral(t, infix.Left, tt.leftValue) {
			return
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
		program := parser.Parse()

		failIfParserHasErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.Expression)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.Expression. got=%T", program.Statements[0])
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

func TestPrefixExpressions(t *testing.T) {
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
		program := parser.Parse()

		failIfParserHasErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.Expression)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.Expression. got=%T", program.Statements[0])
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
		program := parser.Parse()

		failIfParserHasErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.Expression)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.Expression. got=%T", program.Statements[0])
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

func TestListLiteral(t *testing.T) {
	input := `[1, 4, 6]`

	scanner := scanner.New(input)
	tokens := scanner.ScanTokens()
	parser := New(tokens)
	program := parser.Parse()

	failIfParserHasErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.Expression)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.Expression. got=%T", program.Statements[0])
	}

	list, ok := statement.Expression.(*ast.List)

	if !ok {
		t.Fatalf("statement is not ast.List. got=%T", statement.Expression)
	}

	if len(list.Elements) != 3 {
		t.Fatalf("len(list.Elements) is not 3. got=%d", len(list.Elements))
	}

	isNumberLiteral(t, list.Elements[0], 1)
	isNumberLiteral(t, list.Elements[1], 4)
	isNumberLiteral(t, list.Elements[2], 6)
}

func TestIndexExpressions(t *testing.T) {
	input := `example[1 + 1]`

	scanner := scanner.New(input)
	tokens := scanner.ScanTokens()
	parser := New(tokens)
	program := parser.Parse()

	failIfParserHasErrors(t, parser)

	statement, ok := program.Statements[0].(*ast.Expression)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.Expression. got=%T", program.Statements[0])
	}

	index, ok := statement.Expression.(*ast.Index)

	if !ok {
		t.Fatalf("statement is not ast.Index. got=%T", statement.Expression)
	}

	if !isIdentifier(t, index.Left, "example") {
		return
	}
}

// =============================================================================
// Helper methods

func failIfParserHasErrors(t *testing.T, parser *Parser) {
	errors := parser.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, message := range errors {
		t.Errorf("parser error: %q", message)
	}

	t.FailNow()
}

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

func isIdentifier(t *testing.T, expression ast.ExpressionNode, value string) bool {
	identifier, ok := expression.(*ast.Identifier)

	if !ok {
		t.Errorf("expression is not ast.Identifier. got=%T", expression)
	}

	if identifier.Value != value {
		t.Errorf("identifier.Value is not %s. got=%s", value, identifier.Value)
	}

	return true
}
