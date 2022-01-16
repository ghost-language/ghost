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
		parser := New(scanner)
		program := parser.Parse()

		failIfParserHasErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		assign, ok := program.Statements[0].(*ast.Assign)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.Assign. got=%T", program.Statements[0])
		}

		if assign.Name.(*ast.Identifier).Value != tt.identifier {
			t.Fatalf("assign.Name is not '%s'. got=%s", tt.identifier, assign.Name.(*ast.Identifier).Value)
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
		parser := New(scanner)
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

func TestForExpression(t *testing.T) {
	input := `for (x := 0; x < 10; x := x + 1) { true }`

	scanner := scanner.New(input)
	parser := New(scanner)
	program := parser.Parse()

	failIfParserHasErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 Statement. got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.Expression)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.Expression. got=%T", program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.For)

	if !ok {
		t.Fatalf("statement.Expression is not ast.For. got=%T", statement.Expression)
	}

	if !isIdentifier(t, expression.Identifier, "x") {
		return
	}

	if _, ok = expression.Initializer.(*ast.Assign); !ok {
		t.Fatalf("expression.Initializer is not ast.Assign. got=%T", expression.Initializer)
	}

	if _, ok = expression.Increment.(*ast.Assign); !ok {
		t.Fatalf("expression.Increment is not ast.Assign. got=%T", expression.Increment)
	}

	if _, ok = expression.Block.Statements[0].(*ast.Expression); !ok {
		t.Fatalf("expression.Block.Statements[0] is not ast.Expression. got=%T", expression.Block.Statements[0])
	}
}

func TestForInListExpression(t *testing.T) {
	input := `for (x in bar) { true }`

	scanner := scanner.New(input)
	parser := New(scanner)
	program := parser.Parse()

	failIfParserHasErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 Statement. got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.Expression)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.Expression. got=%T", program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.ForIn)

	if !ok {
		t.Fatalf("statement.Expression is not ast.ForIn. got=%T", statement.Expression)
	}

	if !isIdentifier(t, expression.Value, "x") {
		return
	}

	if _, ok = expression.Block.Statements[0].(*ast.Expression); !ok {
		t.Fatalf("expression.Block.Statements[0] is not ast.Expression. got=%T", expression.Block.Statements[0])
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
		parser := New(scanner)
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
	input := `if (x < y) { x }`

	scanner := scanner.New(input)
	parser := New(scanner)
	program := parser.Parse()

	failIfParserHasErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.Expression)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.Expression. got=%T", program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.If)

	if !ok {
		t.Fatalf("statement is not ast.If. got=%T", statement.Expression)
	}

	if !isInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. got=%d", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.Expression)

	if !ok {
		t.Fatalf("Consequence.Statements[0] is not ast.Expression. got=%T", expression.Consequence.Statements[0])
	}

	if !isIdentifier(t, consequence.Expression, "x") {
		return
	}

	if expression.Alternative != nil {
		t.Errorf("expression.Alternative was not nil. got=%+v", expression.Alternative)
	}
}

func TestIfElseExpressions(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	scanner := scanner.New(input)
	parser := New(scanner)
	program := parser.Parse()

	failIfParserHasErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.Expression)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.Expression. got=%T", program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.If)

	if !ok {
		t.Fatalf("statement is not ast.If. got=%T", statement.Expression)
	}

	if !isInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. got=%d", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.Expression)

	if !ok {
		t.Fatalf("Consequence.Statements[0] is not ast.Expression. got=%T", expression.Consequence.Statements[0])
	}

	if !isIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(expression.Alternative.Statements) != 1 {
		t.Errorf("expression.Alternative is not 1 statement. got=%d", len(expression.Alternative.Statements))
	}

	alternative, ok := expression.Alternative.Statements[0].(*ast.Expression)

	if !ok {
		t.Fatalf("Alternative.Statements[0] is not ast.Expression. got=%T", expression.Alternative.Statements[0])
	}

	if !isIdentifier(t, alternative.Expression, "y") {
		return
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
		{"5 % 5", 5, "%", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"1 .. 10", 1, "..", 10},
	}

	for _, tt := range tests {
		scanner := scanner.New(tt.input)
		parser := New(scanner)
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
		parser := New(scanner)
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
		parser := New(scanner)
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
		parser := New(scanner)
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
	parser := New(scanner)
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
	parser := New(scanner)
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

func TestMapLiteralsWithStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	scanner := scanner.New(input)
	parser := New(scanner)
	program := parser.Parse()

	failIfParserHasErrors(t, parser)

	statement, ok := program.Statements[0].(*ast.Expression)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.Expression. got=%T", program.Statements[0])
	}

	mapLiteral, ok := statement.Expression.(*ast.Map)

	if !ok {
		t.Fatalf("statement is not ast.Map. got=%T", statement.Expression)
	}

	if len(mapLiteral.Pairs) != 3 {
		t.Fatalf("map.Pairs has wrong length. got=%d", len(mapLiteral.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range mapLiteral.Pairs {
		literal, ok := key.(*ast.String)

		if !ok {
			t.Errorf("key is not ast.String. got=%T", key)
		}

		expectedValue := expected[literal.Value]

		isNumberLiteral(t, value, expectedValue)
	}
}

func TestEmptyMapLiterals(t *testing.T) {
	input := `{}`

	scanner := scanner.New(input)
	parser := New(scanner)
	program := parser.Parse()

	failIfParserHasErrors(t, parser)

	statement, ok := program.Statements[0].(*ast.Expression)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.Expression. got=%T", program.Statements[0])
	}

	mapLiteral, ok := statement.Expression.(*ast.Map)

	if !ok {
		t.Fatalf("statement is not ast.Map. got=%T", statement.Expression)
	}

	if len(mapLiteral.Pairs) != 0 {
		t.Fatalf("map.Pairs has wrong length. got=%d", len(mapLiteral.Pairs))
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
		return 5
		return 10
		return 3.14
	`

	scanner := scanner.New(input)
	parser := New(scanner)
	program := parser.Parse()

	failIfParserHasErrors(t, parser)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.Return)

		if !ok {
			t.Fatalf("statement is not ast.Return. got=%T", statement)
		}

		if returnStatement.Token.Lexeme != "return" {
			t.Fatalf("returnStatement.Token.Lexeme is not 'return. got=%q", returnStatement.Token.Lexeme)
		}
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

func isInfixExpression(t *testing.T, expression ast.ExpressionNode, left interface{}, operator string, right interface{}) bool {
	operatorExpression, ok := expression.(*ast.Infix)

	if !ok {
		t.Errorf("expression is not ast.Infix. got=%T(%s)", expression, expression)
		return false
	}

	if !isLiteral(t, operatorExpression.Left, left) {
		return false
	}

	if operatorExpression.Operator != operator {
		t.Errorf("expression.Operator is not '%s'. got=%q", operator, operatorExpression.Operator)
		return false
	}

	if !isLiteral(t, operatorExpression.Right, right) {
		return false
	}

	return true
}

func isLiteral(t *testing.T, expression ast.ExpressionNode, expected interface{}) bool {
	switch value := expected.(type) {
	case int:
		return isNumberLiteral(t, expression, int64(value))
	case int64:
		return isNumberLiteral(t, expression, int64(value))
	case float64:
		return isNumberLiteral(t, expression, int64(value))
	case string:
		return isIdentifier(t, expression, value)
	}

	t.Errorf("type of expression is not a literal. got=%T", expression)

	return false
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
