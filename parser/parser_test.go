package parser

import (
	"strings"
	"testing"
	"time"

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
		scanner := scanner.New(tt.input, "test.gs")
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

func TestFunctionRestParameter(t *testing.T) {
	input := `function sum(a, ...rest) { return a }`

	scanner := scanner.New(input, "test.gs")
	parser := New(scanner)
	program := parser.Parse()

	failIfParserHasErrors(t, parser)

	statement, ok := program.Statements[0].(*ast.Expression)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.Expression. got=%T", program.Statements[0])
	}

	function, ok := statement.Expression.(*ast.Function)

	if !ok {
		t.Fatalf("statement is not ast.Function. got=%T", statement.Expression)
	}

	if !function.Rest {
		t.Fatal("function.Rest is false, expected true")
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function.Parameters has wrong length. got=%d", len(function.Parameters))
	}

	if function.Parameters[0].Value != "a" || function.Parameters[1].Value != "rest" {
		t.Fatalf("wrong parameter names. got=%v", function.Parameters)
	}
}

func TestSpreadExpression(t *testing.T) {
	tests := []struct {
		input string
	}{
		{`f(...a)`},
		{`f(1, ...a, 2)`},
		{`[...a, 1]`},
	}

	for _, tt := range tests {
		scanner := scanner.New(tt.input, "test.gs")
		parser := New(scanner)
		program := parser.Parse()

		failIfParserHasErrors(t, parser)

		statement, ok := program.Statements[0].(*ast.Expression)

		if !ok {
			t.Fatalf("%q: program.Statements[0] is not ast.Expression. got=%T", tt.input, program.Statements[0])
		}

		var spread *ast.Spread

		switch expression := statement.Expression.(type) {
		case *ast.Call:
			for _, argument := range expression.Arguments {
				if candidate, ok := argument.(*ast.Spread); ok {
					spread = candidate
				}
			}
		case *ast.List:
			for _, element := range expression.Elements {
				if candidate, ok := element.(*ast.Spread); ok {
					spread = candidate
				}
			}
		default:
			t.Fatalf("%q: unexpected expression type %T", tt.input, statement.Expression)
		}

		if spread == nil {
			t.Fatalf("%q: no ast.Spread found", tt.input)
		}

		if !isIdentifier(t, spread.Value, "a") {
			return
		}
	}
}

func TestListPatternAssignment(t *testing.T) {
	input := `[a, b] = list`

	scanner := scanner.New(input, "test.gs")
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

	pattern, ok := assign.Name.(*ast.ListPattern)

	if !ok {
		t.Fatalf("assign.Name is not ast.ListPattern. got=%T", assign.Name)
	}

	if len(pattern.Targets) != 2 {
		t.Fatalf("pattern.Targets has wrong length. got=%d", len(pattern.Targets))
	}

	if pattern.Targets[0].Value != "a" || pattern.Targets[1].Value != "b" {
		t.Fatalf("wrong pattern targets. got=%v", pattern.Targets)
	}

	if !isIdentifier(t, assign.Value, "list") {
		return
	}
}

func TestMapPatternAssignment(t *testing.T) {
	input := `{x, y: a} = source`

	scanner := scanner.New(input, "test.gs")
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

	pattern, ok := assign.Name.(*ast.MapPattern)

	if !ok {
		t.Fatalf("assign.Name is not ast.MapPattern. got=%T", assign.Name)
	}

	if len(pattern.Pairs) != 2 {
		t.Fatalf("pattern.Pairs has wrong length. got=%d", len(pattern.Pairs))
	}

	byName := map[string]ast.MapPatternPair{}

	for _, pair := range pattern.Pairs {
		byName[pair.Source.Value] = pair
	}

	shorthand, ok := byName["x"]

	if !ok || shorthand.Target.Value != "x" {
		t.Fatalf("shorthand pair for %q not bound correctly. got=%v", "x", byName)
	}

	renamed, ok := byName["y"]

	if !ok || renamed.Target.Value != "a" {
		t.Fatalf("renamed pair for %q not bound correctly. got=%v", "y", byName)
	}

	if !isIdentifier(t, assign.Value, "source") {
		return
	}
}

// A list/map literal that isn't followed by `=` is still an ordinary
// expression statement - the pattern grammar and the literal grammar only
// diverge in meaning, not in what parses.
func TestListAndMapLiteralsAreNotMistakenForPatterns(t *testing.T) {
	tests := []struct {
		input string
	}{
		{`[1, 2, 3]`},
		{`{foo: 1, bar: 2}`},
	}

	for _, tt := range tests {
		scanner := scanner.New(tt.input, "test.gs")
		parser := New(scanner)
		program := parser.Parse()

		failIfParserHasErrors(t, parser)

		if _, ok := program.Statements[0].(*ast.Expression); !ok {
			t.Fatalf("%q: program.Statements[0] is not ast.Expression. got=%T", tt.input, program.Statements[0])
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
		scanner := scanner.New(tt.input, "test.gs")
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

	scanner := scanner.New(input, "test.gs")
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

	scanner := scanner.New(input, "test.gs")
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
		scanner := scanner.New(tt.input, "test.gs")
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

	scanner := scanner.New(input, "test.gs")
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

	scanner := scanner.New(input, "test.gs")
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
		scanner := scanner.New(tt.input, "test.gs")
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
		scanner := scanner.New(tt.input, "test.gs")
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
		scanner := scanner.New(tt.input, "test.gs")
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
		left     interface{}
	}{
		{"index++", "++", "identifier"},
		{"index--", "--", "identifier"},
		{"this.score++", "++", "property"},
		{"list[0]--", "--", "index"},
	}

	for _, tt := range tests {
		scanner := scanner.New(tt.input, "test.gs")
		parser := New(scanner)
		program := parser.Parse()

		failIfParserHasErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("%s: program.Statements does not contain 1 statement. got=%d", tt.input, len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.Expression)

		if !ok {
			t.Fatalf("%s: program.Statements[0] is not ast.Expression. got=%T", tt.input, program.Statements[0])
		}

		postfix, ok := statement.Expression.(*ast.Postfix)

		if !ok {
			t.Fatalf("%s: statement is not ast.Postfix. got=%T", tt.input, statement.Expression)
		}

		if postfix.Operator.String() != tt.operator {
			t.Fatalf("%s: postfix.Operator is not '%s'. got=%s", tt.input, tt.operator, postfix.Operator)
		}

		switch tt.left {
		case "identifier":
			if _, ok := postfix.Left.(*ast.Identifier); !ok {
				t.Fatalf("%s: postfix.Left is not ast.Identifier. got=%T", tt.input, postfix.Left)
			}
		case "property":
			if _, ok := postfix.Left.(*ast.Property); !ok {
				t.Fatalf("%s: postfix.Left is not ast.Property. got=%T", tt.input, postfix.Left)
			}
		case "index":
			if _, ok := postfix.Left.(*ast.Index); !ok {
				t.Fatalf("%s: postfix.Left is not ast.Index. got=%T", tt.input, postfix.Left)
			}
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
		scanner := scanner.New(tt.input, "test.gs")
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

func TestTemplateStringLiteral(t *testing.T) {
	input := "`count: ${1 + 2} done`"

	scanner := scanner.New(input, "test.gs")
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

	template, ok := statement.Expression.(*ast.TemplateString)

	if !ok {
		t.Fatalf("statement is not ast.TemplateString. got=%T", statement.Expression)
	}

	expectedChunks := []string{"count: ", " done"}

	if len(template.Chunks) != len(expectedChunks) {
		t.Fatalf("template.Chunks has wrong length. expected=%d, got=%d", len(expectedChunks), len(template.Chunks))
	}

	for index, chunk := range expectedChunks {
		if template.Chunks[index] != chunk {
			t.Errorf("template.Chunks[%d] is not %q. got=%q", index, chunk, template.Chunks[index])
		}
	}

	if len(template.Expressions) != 1 {
		t.Fatalf("template.Expressions has wrong length. expected=1, got=%d", len(template.Expressions))
	}

	isInfixExpression(t, template.Expressions[0], 1, "+", 2)
}

func TestTemplateStringLiteralWithNoInterpolation(t *testing.T) {
	input := "`hello`"

	scanner := scanner.New(input, "test.gs")
	parser := New(scanner)
	program := parser.Parse()

	failIfParserHasErrors(t, parser)

	statement := program.Statements[0].(*ast.Expression)
	template, ok := statement.Expression.(*ast.TemplateString)

	if !ok {
		t.Fatalf("statement is not ast.TemplateString. got=%T", statement.Expression)
	}

	if len(template.Expressions) != 0 {
		t.Fatalf("template.Expressions should be empty. got=%d", len(template.Expressions))
	}

	if len(template.Chunks) != 1 || template.Chunks[0] != "hello" {
		t.Fatalf("template.Chunks is wrong. got=%v", template.Chunks)
	}
}

func TestListLiteral(t *testing.T) {
	input := `[1, 4, 6]`

	scanner := scanner.New(input, "test.gs")
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

	scanner := scanner.New(input, "test.gs")
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

	scanner := scanner.New(input, "test.gs")
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

	scanner := scanner.New(input, "test.gs")
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

	scanner := scanner.New(input, "test.gs")
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

	scanner := scanner.New(input, "test.gs")
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

func TestMapLiteralsWithShorthandKeys(t *testing.T) {
	input := `{foo, bar, baz: 3}`

	scanner := scanner.New(input, "test.gs")
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

	for key, value := range mapLiteral.Pairs {
		keyIdentifier, ok := key.(*ast.Identifier)

		if !ok {
			t.Fatalf("key is not ast.Identifier. got=%T", key)
		}

		if keyIdentifier.Value == "baz" {
			isNumberLiteral(t, value, 3)
			continue
		}

		valueIdentifier, ok := value.(*ast.Identifier)

		if !ok {
			t.Fatalf("value for shorthand key %q is not ast.Identifier. got=%T", keyIdentifier.Value, value)
		}

		if valueIdentifier.Value != keyIdentifier.Value {
			t.Errorf("shorthand value should reference %q, got=%q", keyIdentifier.Value, valueIdentifier.Value)
		}
	}
}

func TestEmptyMapLiterals(t *testing.T) {
	input := `{}`

	scanner := scanner.New(input, "test.gs")
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

	scanner := scanner.New(input, "test.gs")
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

	scanner := scanner.New(input, "test.gs")
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

	scanner := scanner.New(input, "test.gs")
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

	scanner := scanner.New(input, "test.gs")
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

	scanner := scanner.New(input, "test.gs")
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
		scanner := scanner.New(tt.input, "test.gs")
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
	scanner := scanner.New(`new Person().greet()`, "test.gs")
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

	scanner := scanner.New(input, "test.gs")
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
	scanner := scanner.New(`super.greet()`, "test.gs")
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
		{`Person.new()`, "test.gs:1:8: syntax error: `new` is not a method"},
		{`x = ,`, "test.gs:1:5: syntax error: `,` cannot start an expression"},
		{`[a, b+1] = list`, "test.gs:1:1: syntax error: a list pattern can only bind plain names, not a full expression"},
		{`{x: 1+1} = source`, "test.gs:1:1: syntax error: a map pattern can only bind to a plain name, not a full expression"},
		{`function f(...rest, a) { return a }`, "test.gs:1:21: syntax error: a rest parameter has to be the last one"},
		{`function f(...rest = 1) { return rest }`, "test.gs:1:20: syntax error: a rest parameter cannot have a default value"},
	}

	for _, tt := range tests {
		scanner := scanner.New(tt.input, "test.gs")
		parser := New(scanner)
		parser.Parse()

		if len(parser.Errors()) == 0 {
			t.Fatalf("expected a parser error for %q, got none", tt.input)
		}

		if parser.Errors()[0].String() != tt.expected {
			t.Fatalf("wrong parser error. got=%s, expected=%s", parser.Errors()[0], tt.expected)
		}
	}
}

// A file with two mistakes in it should report two mistakes. Without recovery
// the parser reads everything after the first one out of context, and one typo
// arrives as a wall of errors that say nothing about each other.
func TestParserRecoversAndReportsEveryError(t *testing.T) {
	source := "x = 1\ny = )\nz = 2\n\nfunction ok() { return 3 }\n\nw = ,\n"

	parser := New(scanner.New(source, "test.gs"))
	parser.Parse()

	errors := parser.Errors()

	if len(errors) != 2 {
		t.Fatalf("got %d errors, expected 2: %v", len(errors), errors)
	}

	expected := []string{
		"test.gs:2:5: syntax error: `)` cannot start an expression",
		"test.gs:7:5: syntax error: `,` cannot start an expression",
	}

	for index, want := range expected {
		if errors[index].String() != want {
			t.Errorf("error %d: got=%q, expected=%q", index, errors[index].String(), want)
		}
	}
}

// Which pass noticed a problem is Ghost's business, not the reader's, so
// lexical and grammatical errors come back as one list in source order.
func TestParserFoldsInScannerFaults(t *testing.T) {
	parser := New(scanner.New("a = @\nb = )\n", "test.gs"))
	parser.Parse()

	errors := parser.Errors()

	if len(errors) < 2 {
		t.Fatalf("got %d errors, expected at least 2: %v", len(errors), errors)
	}

	if errors[0].Position.Line > errors[1].Position.Line {
		t.Errorf("errors are out of source order: %v", errors)
	}

	if errors[0].String() != "test.gs:1:5: syntax error: unexpected character `@`" {
		t.Errorf("got=%q", errors[0].String())
	}
}

// A number the scanner accepted but that Ghost cannot hold is a syntax error
// with a position, not a line logged to the terminal from inside the parser.
func TestParserReportsUnreadableNumbers(t *testing.T) {
	tests := []struct {
		source   string
		expected string
	}{
		{"x = 99999999999999999999999999", "test.gs:1:5: syntax error: `99999999999999999999999999` is not a valid number"},
		{"x = 1e", "test.gs:1:5: syntax error: `1e` is not a valid number"},
	}

	for _, test := range tests {
		parser := New(scanner.New(test.source, "test.gs"))
		parser.Parse()

		errors := parser.Errors()

		if len(errors) == 0 {
			t.Fatalf("expected an error for %q", test.source)
		}

		if errors[0].String() != test.expected {
			t.Errorf("got=%q, expected=%q", errors[0].String(), test.expected)
		}
	}
}

// A closing bracket that starts an expression is nearly always an unclosed
// opener somewhere above it, and saying so is more use than naming the token.
func TestImportStatements(t *testing.T) {
	t.Run("bare import has no alias", func(t *testing.T) {
		program := New(scanner.New(`import "math"`, "test.gs")).Parse()

		statement := program.Statements[0].(*ast.Expression)
		imp, ok := statement.Expression.(*ast.Import)

		if !ok {
			t.Fatalf("expected ast.Import, got %T", statement.Expression)
		}

		if imp.Path.Value != "math" {
			t.Fatalf("expected path %q, got %q", "math", imp.Path.Value)
		}

		if imp.Alias != nil {
			t.Fatalf("expected no alias, got %q", imp.Alias.Value)
		}
	})

	t.Run("bare import as alias", func(t *testing.T) {
		program := New(scanner.New(`import "math" as m`, "test.gs")).Parse()

		statement := program.Statements[0].(*ast.Expression)
		imp := statement.Expression.(*ast.Import)

		if imp.Alias == nil || imp.Alias.Value != "m" {
			t.Fatalf("expected alias %q, got %v", "m", imp.Alias)
		}
	})

	t.Run("named import without braces still works", func(t *testing.T) {
		program := New(scanner.New(`import pi, e from "math"`, "test.gs")).Parse()

		statement := program.Statements[0].(*ast.Expression)
		importFrom := statement.Expression.(*ast.ImportFrom)

		if len(importFrom.Identifiers) != 2 {
			t.Fatalf("expected 2 identifiers, got %d", len(importFrom.Identifiers))
		}
	})

	t.Run("named import with braces", func(t *testing.T) {
		program := New(scanner.New(`import { pi, e } from "math"`, "test.gs")).Parse()

		statement := program.Statements[0].(*ast.Expression)
		importFrom := statement.Expression.(*ast.ImportFrom)

		if len(importFrom.Identifiers) != 2 {
			t.Fatalf("expected 2 identifiers, got %d", len(importFrom.Identifiers))
		}

		if importFrom.Path.Value != "math" {
			t.Fatalf("expected path %q, got %q", "math", importFrom.Path.Value)
		}
	})

	t.Run("named import with braces and alias", func(t *testing.T) {
		program := New(scanner.New(`import { pi as p } from "math"`, "test.gs")).Parse()

		statement := program.Statements[0].(*ast.Expression)
		importFrom := statement.Expression.(*ast.ImportFrom)

		identifier, ok := importFrom.Identifiers["p"]

		if !ok {
			t.Fatalf("expected alias %q to be bound, got %v", "p", importFrom.Identifiers)
		}

		if identifier.Value != "pi" {
			t.Fatalf("expected alias %q to name %q, got %q", "p", "pi", identifier.Value)
		}
	})

	t.Run("everything import with braces", func(t *testing.T) {
		program := New(scanner.New(`import { * } from "math"`, "test.gs")).Parse()

		statement := program.Statements[0].(*ast.Expression)
		importFrom := statement.Expression.(*ast.ImportFrom)

		if !importFrom.Everything {
			t.Fatal("expected Everything to be true")
		}
	})

	t.Run("unclosed brace reports an error", func(t *testing.T) {
		parser := New(scanner.New(`import { pi, e from "math"`, "test.gs"))
		parser.Parse()

		if len(parser.Errors()) == 0 {
			t.Fatal("expected an error for an unclosed brace")
		}
	})

	t.Run("combined module and named import (JS-style)", func(t *testing.T) {
		program := New(scanner.New(`import image, { Spritesheet } from "lumen:image"`, "test.gs")).Parse()

		statement := program.Statements[0].(*ast.Expression)
		imp, ok := statement.Expression.(*ast.Import)

		if !ok {
			t.Fatalf("expected ast.Import, got %T", statement.Expression)
		}

		if imp.Path.Value != "lumen:image" {
			t.Fatalf("expected path %q, got %q", "lumen:image", imp.Path.Value)
		}

		if imp.Alias == nil || imp.Alias.Value != "image" {
			t.Fatalf("expected the module bound to %q, got %v", "image", imp.Alias)
		}

		identifier, ok := imp.Identifiers["Spritesheet"]

		if !ok {
			t.Fatalf("expected %q to be bound, got %v", "Spritesheet", imp.Identifiers)
		}

		if identifier.Value != "Spritesheet" {
			t.Fatalf("expected identifier %q, got %q", "Spritesheet", identifier.Value)
		}
	})

	t.Run("combined module and named import with a named alias", func(t *testing.T) {
		program := New(scanner.New(`import image, { Spritesheet as Sheet } from "lumen:image"`, "test.gs")).Parse()

		statement := program.Statements[0].(*ast.Expression)
		imp := statement.Expression.(*ast.Import)

		if imp.Alias == nil || imp.Alias.Value != "image" {
			t.Fatalf("expected the module bound to %q, got %v", "image", imp.Alias)
		}

		identifier, ok := imp.Identifiers["Sheet"]

		if !ok {
			t.Fatalf("expected alias %q to be bound, got %v", "Sheet", imp.Identifiers)
		}

		if identifier.Value != "Spritesheet" {
			t.Fatalf("expected alias %q to name %q, got %q", "Sheet", "Spritesheet", identifier.Value)
		}
	})

	t.Run("combined module and everything import", func(t *testing.T) {
		program := New(scanner.New(`import image, { * } from "lumen:image"`, "test.gs")).Parse()

		statement := program.Statements[0].(*ast.Expression)
		imp := statement.Expression.(*ast.Import)

		if !imp.Everything {
			t.Fatal("expected Everything to be true")
		}
	})

	t.Run("combined form requires braces around the named list", func(t *testing.T) {
		parser := New(scanner.New(`import image, Spritesheet from "lumen:image"`, "test.gs"))
		program := parser.Parse()

		// `image, Spritesheet` with no brace isn't the combined form at all —
		// it's an ordinary unbraced named-import list of two names, so this
		// parses fine but pulls two names (`image`, `Spritesheet`) rather than
		// binding the module.
		if len(parser.Errors()) != 0 {
			t.Fatalf("expected no error, got %v", parser.Errors())
		}

		statement := program.Statements[0].(*ast.Expression)
		importFrom, ok := statement.Expression.(*ast.ImportFrom)

		if !ok {
			t.Fatalf("expected ast.ImportFrom, got %T", statement.Expression)
		}

		if len(importFrom.Identifiers) != 2 {
			t.Fatalf("expected 2 identifiers, got %d", len(importFrom.Identifiers))
		}
	})

	t.Run("unbraced named list still works alongside the combined form", func(t *testing.T) {
		program := New(scanner.New(`import pi, e, tau from "math"`, "test.gs")).Parse()

		statement := program.Statements[0].(*ast.Expression)
		importFrom := statement.Expression.(*ast.ImportFrom)

		if len(importFrom.Identifiers) != 3 {
			t.Fatalf("expected 3 identifiers, got %d", len(importFrom.Identifiers))
		}
	})
}

func TestParserSuggestsAMissingOpener(t *testing.T) {
	parser := New(scanner.New("x = )", "test.gs"))
	parser.Parse()

	errors := parser.Errors()

	if len(errors) == 0 {
		t.Fatal("expected an error")
	}

	if errors[0].Help == "" {
		t.Errorf("expected help on %q", errors[0])
	}
}

// A parser that never finishes is a worse failure than any error message: the
// program does not run and does not say why. These are the two loops that used
// to spin at the end of the file rather than reporting what was missing.
func TestParsingAlwaysTerminates(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected string
	}{
		{"an import written backwards", `from "lib" import double`, "expected `from` after the names being imported"},
		{"an unclosed parameter list", `function greet(name, greeting`, "expected `)` to close the parameter list"},
		{"a parameter that is not a name", `function greet(1) { }`, "expected a parameter name, found `1`"},
		{"an import of something unnamed", `import 5 from "lib"`, "expected a name to import, found `5`"},
		{"an unclosed map", `x = {"a": 1`, "expected"},
		{"an unclosed switch", `switch (x) { case 1: {`, ""},
		{"an unclosed block", `function greet() {`, ""},
		{"an unclosed class", `class Point {`, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			finished := make(chan []string, 1)

			go func() {
				parser := New(scanner.New(test.source, "test.gs"))
				parser.Parse()

				messages := make([]string, 0, len(parser.Errors()))

				for _, raised := range parser.Errors() {
					messages = append(messages, raised.Message)
				}

				finished <- messages
			}()

			select {
			case messages := <-finished:
				if test.expected == "" {
					return
				}

				for _, message := range messages {
					if strings.Contains(message, test.expected) {
						return
					}
				}

				t.Errorf("no error mentioned %q, got %v", test.expected, messages)
			case <-time.After(5 * time.Second):
				t.Fatalf("parsing %q did not finish", test.source)
			}
		})
	}
}
