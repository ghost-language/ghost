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
		{`a = 5`, "a", 5},
	}

	for _, tt := range tests {
		scanner := scanner.New(tt.input, "test.ghost")
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
		scanner := scanner.New(tt.input, "test.ghost")
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
	input := `for (x = 0; x < 10; x = x + 1) { true }`

	scanner := scanner.New(input, "test.ghost")
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

	scanner := scanner.New(input, "test.ghost")
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
		scanner := scanner.New(tt.input, "test.ghost")
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

	scanner := scanner.New(input, "test.ghost")
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

	scanner := scanner.New(input, "test.ghost")
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
		scanner := scanner.New(tt.input, "test.ghost")
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

		if infix.Operator.String() != tt.operator {
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
		input      string
		isFloat    bool
		intValue   int64
		floatValue float64
	}{
		{"5", false, 5, 0},
		{"3.14", true, 0, 3.14},
		{"5e10", true, 0, 5e10},
	}

	for _, tt := range tests {
		scanner := scanner.New(tt.input, "test.ghost")
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

		if number.IsFloat != tt.isFloat {
			t.Fatalf("number.IsFloat is not '%t'. got=%t", tt.isFloat, number.IsFloat)
		}

		if tt.isFloat {
			if number.FloatValue != tt.floatValue {
				t.Fatalf("number.FloatValue is not '%g'. got=%g", tt.floatValue, number.FloatValue)
			}
		} else {
			if number.IntValue != tt.intValue {
				t.Fatalf("number.IntValue is not '%d'. got=%d", tt.intValue, number.IntValue)
			}
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
		scanner := scanner.New(tt.input, "test.ghost")
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

		if prefix.Operator.String() != tt.operator {
			t.Fatalf("prefix.Operator is not '%s'. got=%s", tt.operator, prefix.Operator)
		}

		if !isNumberLiteral(t, prefix.Right, tt.number) {
			return
		}
	}
}

func TestPostfixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
	}{
		{"index++", "++"},
		{"index--", "--"},
	}

	for _, tt := range tests {
		scanner := scanner.New(tt.input, "test.ghost")
		parser := New(scanner)
		program := parser.Parse()

		failIfParserHasErrors(t, parser)

		if len(program.Statements) != 2 {
			t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
		}

		statement, ok := program.Statements[1].(*ast.Expression)

		if !ok {
			t.Fatalf("program.Statements[1] is not ast.Expression. got=%T", program.Statements[0])
		}

		postfix, ok := statement.Expression.(*ast.Postfix)

		if !ok {
			t.Fatalf("statement is not ast.Postfix. got=%T", statement.Expression)
		}

		if postfix.Operator.String() != tt.operator {
			t.Fatalf("postfix.Operator is not '%s'. got=%s", tt.operator, postfix.Operator)
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
		scanner := scanner.New(tt.input, "test.ghost")
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

	scanner := scanner.New(input, "test.ghost")
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

	scanner := scanner.New(input, "test.ghost")
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

	scanner := scanner.New(input, "test.ghost")
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

func TestMapLiteralsWithBooleanKeys(t *testing.T) {
	input := `{true: 1, false: 2}`

	scanner := scanner.New(input, "test.ghost")
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

	if len(mapLiteral.Pairs) != 2 {
		t.Fatalf("map.Pairs has wrong length. got=%d", len(mapLiteral.Pairs))
	}

	expected := map[bool]int64{
		true:  1,
		false: 2,
	}

	for key, value := range mapLiteral.Pairs {
		boolean, ok := key.(*ast.Boolean)

		if !ok {
			t.Errorf("key is not ast.Boolean. got=%T", key)
		}

		expectedValue := expected[boolean.Value]

		isNumberLiteral(t, value, expectedValue)
	}
}

func TestMapLiteralsWithIntegerKeys(t *testing.T) {
	input := `{1: 1, 2: 2, 3: 3}`

	scanner := scanner.New(input, "test.ghost")
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

	expected := map[int64]int64{
		1: 1,
		2: 2,
		3: 3,
	}

	for key, value := range mapLiteral.Pairs {
		number, ok := key.(*ast.Number)

		if !ok {
			t.Errorf("key is not ast.Number. got=%T", key)
		}

		expectedValue := expected[number.IntValue]

		isNumberLiteral(t, value, expectedValue)
	}
}

func TestMapLiteralsWithVariableKeys(t *testing.T) {
	input := `{foo: 1, bar: 2, baz: 3}`

	scanner := scanner.New(input, "test.ghost")
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
		"foo": 1,
		"bar": 2,
		"baz": 3,
	}

	for key, value := range mapLiteral.Pairs {
		identifier, ok := key.(*ast.Identifier)

		if !ok {
			t.Errorf("key is not ast.Identifier. got=%T", key)
		}

		expectedValue := expected[identifier.Value]

		isNumberLiteral(t, value, expectedValue)
	}
}

func TestEmptyMapLiterals(t *testing.T) {
	input := `{}`

	scanner := scanner.New(input, "test.ghost")
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

	scanner := scanner.New(input, "test.ghost")
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

func TestSwitchStatements(t *testing.T) {
	input := `switch (value) {
		case 1 {
			print('one')
		}
		case 2 {
			print('two')
		}
	}`

	scanner := scanner.New(input, "test.ghost")
	parser := New(scanner)
	program := parser.Parse()

	failIfParserHasErrors(t, parser)

	statement, ok := program.Statements[0].(*ast.Expression)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.Expression. got=%T", program.Statements[0])
	}

	switchStatement, ok := statement.Expression.(*ast.Switch)

	if !ok {
		t.Fatalf("statement is not ast.Switch. got=%T", statement.Expression)
	}

	if len(switchStatement.Cases) != 2 {
		t.Fatalf("switchStatement.Cases has wrong length. got=%d", len(switchStatement.Cases))
	}
}

func TestSwitchStatementsWithDefault(t *testing.T) {
	input := `switch (value) {
		case 1 {
			print('one')
		}
		case 2 {
			print('two')
		}
		default {
			print('default')
		}
	}`

	scanner := scanner.New(input, "test.ghost")
	parser := New(scanner)
	program := parser.Parse()

	failIfParserHasErrors(t, parser)

	statement, ok := program.Statements[0].(*ast.Expression)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.Expression. got=%T", program.Statements[0])
	}

	switchStatement, ok := statement.Expression.(*ast.Switch)

	if !ok {
		t.Fatalf("statement is not ast.Switch. got=%T", statement.Expression)
	}

	if len(switchStatement.Cases) != 3 {
		t.Fatalf("switchStatement.Cases has wrong length. got=%d", len(switchStatement.Cases))
	}
}

func TestSwitchStatementsWithMultipleDefaults(t *testing.T) {
	input := `switch (value) {
		case 1 {
			print('one')
		}
		case 2 {
			print('two')
		}
		case default {
			print('default one')
		}
		default {
			print('default two')
		}
	}`

	scanner := scanner.New(input, "test.ghost")
	parser := New(scanner)
	parser.Parse()

	// Expecting a parser error here for having multiple defaults
	if len(parser.Errors()) != 1 {
		t.Fatalf("parser should have 1 error. got=%d", len(parser.Errors()))
	}
}

func TestTraitExpressions(t *testing.T) {
	input := `trait Foo {
		//
	}`

	scanner := scanner.New(input, "test.ghost")
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

	expression, ok := statement.Expression.(*ast.Trait)

	if !ok {
		t.Fatalf("statement is not ast.Trait. got=%T", statement.Expression)
	}

	if expression.Name.Value != "Foo" {
		t.Fatalf("expression.Name is not 'Foo'. got=%s", expression.Name.Value)
	}

	if len(expression.Body.Statements) != 0 {
		t.Fatalf("expression.Body.Statements does not contain 0 statements. got=%d", len(expression.Body.Statements))
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

	if operatorExpression.Operator.String() != operator {
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

	if number.IntValue != value {
		t.Errorf("number.IntValue is not %d. got=%d", value, number.IntValue)
	}

	return true
}

func TestNewExpression(t *testing.T) {
	tests := []struct {
		input     string
		class     string
		arguments int
	}{
		{`new Person()`, "Person", 0},
		{`new Person`, "Person", 0},
		{`new Person("kai", 5)`, "Person", 2},
	}

	for _, tt := range tests {
		scanner := scanner.New(tt.input, "test.ghost")
		parser := New(scanner)
		program := parser.Parse()

		failIfParserHasErrors(t, parser)

		statement, ok := program.Statements[0].(*ast.Expression)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.Expression. got=%T", program.Statements[0])
		}

		expression, ok := statement.Expression.(*ast.New)

		if !ok {
			t.Fatalf("statement.Expression is not ast.New. got=%T", statement.Expression)
		}

		identifier, ok := expression.Class.(*ast.Identifier)

		if !ok {
			t.Fatalf("expression.Class is not ast.Identifier. got=%T", expression.Class)
		}

		if identifier.Value != tt.class {
			t.Fatalf("expression.Class is not '%s'. got=%s", tt.class, identifier.Value)
		}

		if len(expression.Arguments) != tt.arguments {
			t.Fatalf("expression.Arguments does not contain %d arguments. got=%d", tt.arguments, len(expression.Arguments))
		}
	}
}

// TestNewBindsTighterThanCalls confirms `new Foo().bar()` reads as
// `(new Foo()).bar()` rather than instantiating the result of `Foo().bar()`.
func TestNewBindsTighterThanCalls(t *testing.T) {
	scanner := scanner.New(`new Person().greet()`, "test.ghost")
	parser := New(scanner)
	program := parser.Parse()

	failIfParserHasErrors(t, parser)

	statement := program.Statements[0].(*ast.Expression)

	method, ok := statement.Expression.(*ast.Method)

	if !ok {
		t.Fatalf("statement.Expression is not ast.Method. got=%T", statement.Expression)
	}

	if _, ok := method.Left.(*ast.New); !ok {
		t.Fatalf("method.Left is not ast.New. got=%T", method.Left)
	}
}

func TestClassMethodShorthand(t *testing.T) {
	input := `
	class Person {
		name = "kai"

		constructor(name) {
			this.name = name
		}

		greet(greeting = "hello") {
			return greeting
		}
	}
	`

	scanner := scanner.New(input, "test.ghost")
	parser := New(scanner)
	program := parser.Parse()

	failIfParserHasErrors(t, parser)

	statement := program.Statements[0].(*ast.Expression)
	class := statement.Expression.(*ast.Class)

	if len(class.Body.Statements) != 3 {
		t.Fatalf("class body does not contain 3 members. got=%d", len(class.Body.Statements))
	}

	if _, ok := class.Body.Statements[0].(*ast.Assign); !ok {
		t.Fatalf("first member is not ast.Assign. got=%T", class.Body.Statements[0])
	}

	for index, expected := range []struct {
		name       string
		parameters int
	}{{"constructor", 1}, {"greet", 1}} {
		member, ok := class.Body.Statements[index+1].(*ast.Expression)

		if !ok {
			t.Fatalf("member %d is not ast.Expression. got=%T", index+1, class.Body.Statements[index+1])
		}

		function, ok := member.Expression.(*ast.Function)

		if !ok {
			t.Fatalf("member %d is not ast.Function. got=%T", index+1, member.Expression)
		}

		if function.Name.Value != expected.name {
			t.Fatalf("member %d is not named '%s'. got=%s", index+1, expected.name, function.Name.Value)
		}

		if len(function.Parameters) != expected.parameters {
			t.Fatalf("member %d does not have %d parameters. got=%d", index+1, expected.parameters, len(function.Parameters))
		}
	}
}

func TestSuperExpression(t *testing.T) {
	scanner := scanner.New(`super.greet()`, "test.ghost")
	parser := New(scanner)
	program := parser.Parse()

	failIfParserHasErrors(t, parser)

	statement := program.Statements[0].(*ast.Expression)

	method, ok := statement.Expression.(*ast.Method)

	if !ok {
		t.Fatalf("statement.Expression is not ast.Method. got=%T", statement.Expression)
	}

	if _, ok := method.Left.(*ast.Super); !ok {
		t.Fatalf("method.Left is not ast.Super. got=%T", method.Left)
	}
}

// TestSyntaxErrors covers tokens that cannot begin an expression. These used to
// parse to a nil node that only surfaced later as a crash in the evaluator.
func TestSyntaxErrors(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`Person.new()`, "1:8: syntax error: `new` is not a method; construct instances with `new Person()`"},
		{`x = ,`, "1:5: syntax error: unexpected token `,`"},
	}

	for _, tt := range tests {
		scanner := scanner.New(tt.input, "test.ghost")
		parser := New(scanner)
		parser.Parse()

		if len(parser.Errors()) == 0 {
			t.Fatalf("expected a parser error for %q, got none", tt.input)
		}

		if parser.Errors()[0] != tt.expected {
			t.Fatalf("wrong parser error. got=%s, expected=%s", parser.Errors()[0], tt.expected)
		}
	}
}
