package parser

import (
	"fmt"
	"testing"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/lexer"
)

func TestAssignStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		exptectedValue     interface{}
	}{
		{"x = 5", "x", 5},
		{"y = true", "y", true},
		{"foobar = y", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		statement := program.Statements[0]

		if !testAssignStatement(t, statement, tt.expectedIdentifier) {
			return
		}

		value := statement.(*ast.AssignStatement).Value

		if !testLiteralExpression(t, value, tt.exptectedValue) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue int
	}{
		{"return 5;", 5},
		{"return 10;", 10},
		{"return 993322;", 993322},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		statement := program.Statements[0]
		returnStatement, ok := statement.(*ast.ReturnStatement)

		if !ok {
			t.Fatalf("statement not *ast.ReturnStatement. got=%T", statement)
		}

		if returnStatement.TokenLiteral() != "return" {
			t.Fatalf("returnStatement.TokenLiteral not 'return'. got=%q", returnStatement.TokenLiteral())
		}

		if testLiteralExpression(t, returnStatement.ReturnValue, tt.expectedValue) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	identifier, ok := statement.Expression.(*ast.IdentifierLiteral)

	if !ok {
		t.Fatalf("expression not *ast.Identifier. got=%T", statement.Expression)
	}

	if identifier.Value != "foobar" {
		t.Errorf("identifier.Value not %s. got=%s", "foobar", identifier.Value)
	}

	if identifier.TokenLiteral() != "foobar" {
		t.Errorf("identifier.TokenLiteral not %s. got=%s", "foobar", identifier.TokenLiteral())
	}
}

func TestNumberLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", statement.Expression)
	}

	literal, ok := statement.Expression.(*ast.NumberLiteral)

	if !ok {
		t.Fatalf("expression not *ast.NumberLiteral. got=%T", statement.Expression)
	}

	if literal.Value.IntPart() != 5 {
		t.Errorf("literal.Value not 5. got=%d", literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	statement := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := statement.Expression.(*ast.StringLiteral)

	if !ok {
		t.Fatalf("expression not *ast.StringLiteral. got=%T", statement.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d\n", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		boolean, ok := statement.Expression.(*ast.BooleanLiteral)

		if !ok {
			t.Fatalf("statement is not ast.Boolean. got=%T", statement.Expression)
		}

		if boolean.Value != tt.expectedBoolean {
			t.Fatalf("boolean.Value is not '%t'. got=%t", tt.expectedBoolean, boolean.Value)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d\n", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.IfExpression)

	if !ok {
		t.Fatalf("statement is not ast.IfExpression. got=%T", statement.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. got=%d", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", expression.Consequence.Statements[0])
	}

	if !testIdentifierLiteral(t, consequence.Expression, "x") {
		return
	}

	if expression.Alternative != nil {
		t.Errorf("expression.Alternative.Statements was not nil. got=%+v", expression.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d\n", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.IfExpression)

	if !ok {
		t.Fatalf("statement is not ast.IfExpression. got=%T", statement.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. got=%d", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", expression.Consequence.Statements[0])
	}

	if !testIdentifierLiteral(t, consequence.Expression, "x") {
		return
	}

	if len(expression.Alternative.Statements) != 1 {
		t.Errorf("alternative is not 1 statement. got=%d", len(expression.Alternative.Statements))
	}

	alternative, ok := expression.Alternative.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", expression.Alternative.Statements[0])
	}

	if !testIdentifierLiteral(t, alternative.Expression, "y") {
		return
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d\n", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		expression, ok := statement.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("statement is not ast.PrefixExpression. got=%T", statement.Expression)
		}

		if expression.Operator != tt.operator {
			t.Fatalf("expression.Operator is not '%s'. got=%s", tt.operator, expression.Operator)
		}

		if !testLiteralExpression(t, expression.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
		{"true and true", true, "and", true},
		{"true or false", true, "or", false},
		{"1 .. 10", 1, "..", 10},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d\n", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		if !testInfixExpression(t, statement.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestParsingPostfixExpressions(t *testing.T) {
	postfixTests := []struct {
		input    string
		operator string
	}{
		{"index++", "++"},
		{"index--", "--"},
	}

	for _, tt := range postfixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 2 {
			t.Fatalf("program.Statements does not contain 2 statements. got=%d\n", len(program.Statements))
		}

		statement, ok := program.Statements[1].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[1] is not ast.ExpressionStatement. got=%T", program.Statements[1])
		}

		expression, ok := statement.Expression.(*ast.PostfixExpression)

		if !ok {
			t.Fatalf("statement is not ast.PostfixExpression. got=%T", statement.Expression)
		}

		if expression.Operator != tt.operator {
			t.Fatalf("expression.Operator is not '%s'. got=%s", tt.operator, expression.Operator)
		}
	}
}

func TestParsingImportExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`import("foobar")`,
			`import("foobar")`,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / e + f)",
			"add((((a + b) + ((c * d) / e)) + f))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
		{
			"a % 2",
			"(a % 2)",
		},
		{
			"foo.bar * foo.baz",
			"((foo.bar) * (foo.baz))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `function(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d\n", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := statement.Expression.(*ast.FunctionLiteral)

	if !ok {
		t.Fatalf("statement.Expression is not ast.FunctionLiteral. got=%T", statement.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. expected 2, got=%d\n", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements is not 1 statement. got=%d\n", len(function.Body.Statements))
	}

	bodyStatement, ok := function.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("function body statement is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStatement.Expression, "x", "+", "y")
}

func TestNamedFunctionLiteralParsing(t *testing.T) {
	input := `function test(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d\n", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := statement.Expression.(*ast.FunctionLiteral)

	if !ok {
		t.Fatalf("statement.Expression is not ast.FunctionLiteral. got=%T", statement.Expression)
	}

	if function.Name != "test" {
		t.Fatalf("function literal name wrong. expected 'test', got=%s\n", function.Name)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. expected 2, got=%d\n", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements is not 1 statement. got=%d\n", len(function.Body.Statements))
	}

	bodyStatement, ok := function.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("function body statement is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStatement.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"function() {};", []string{}},
		{"function(x) {};", []string{"x"}},
		{"function(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		statement := program.Statements[0].(*ast.ExpressionStatement)
		function := statement.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. expected %d, got=%d\n", len(tt.expectedParams), len(function.Parameters))
		}

		for i, identifier := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], identifier)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d\n", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.CallExpression)

	if !ok {
		t.Fatalf("statement.Expression is not ast.CallExpression. got=%T", statement.Expression)
	}

	if !testIdentifierLiteral(t, expression.Callable, "add") {
		return
	}

	if len(expression.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(expression.Arguments))
	}

	testLiteralExpression(t, expression.Arguments[0], 1)
	testInfixExpression(t, expression.Arguments[1], 2, "*", 3)
	testInfixExpression(t, expression.Arguments[2], 4, "+", 5)
}

func TestBindExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"x := 5", "x := 5"},
		{"y := true", "y := true"},
		{"foobar := x", "foobar := x"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if program.String() != tt.expected {
			t.Fatalf("bind expression evaluated incorrectly. got=%s, expected=%s", program.String(), tt.expected)
		}
	}
}

func testAssignStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "=" {
		t.Errorf("s.TokenLiteral not '='. got=%q", s.TokenLiteral())
		return false
	}

	assignStatement, ok := s.(*ast.AssignStatement)

	if !ok {
		t.Errorf("s not *ast.AssignStatement. got=%T", s)
		return false
	}

	if assignStatement.Name.Value != name {
		t.Errorf("assignStatement.Name.Value not '%s'. got=%s", name, assignStatement.Name.Value)
		return false
	}

	if assignStatement.Name.TokenLiteral() != name {
		t.Errorf("assignStatement.Name.TokenLiteral() not '%s'. got=%s", name, assignStatement.Name.TokenLiteral())
		return false
	}

	return true
}

func testNumberLiteral(t *testing.T, il ast.Expression, value int64) bool {
	number, ok := il.(*ast.NumberLiteral)

	if !ok {
		t.Errorf("il not *ast.NumberLiteral. got=%T", il)
		return false
	}

	if number.Value.IntPart() != value {
		t.Errorf("number.Value not %d. got=%d", value, number.Value)
		return false
	}

	if number.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("number.TokenLiteral not %d. got=%s", value, number.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, expression ast.Expression, value bool) bool {
	boolean, ok := expression.(*ast.BooleanLiteral)

	if !ok {
		t.Errorf("expression not *ast.Boolean. got=%T", expression)
		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value not %t. got=%t", value, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolean.TokenLiteral not %t. got=%s", value, boolean.TokenLiteral())
		return false
	}

	return true
}

func TestParsingListLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	statement, _ := program.Statements[0].(*ast.ExpressionStatement)
	list, ok := statement.Expression.(*ast.ListLiteral)

	if !ok {
		t.Fatalf("expression not *ast.ListLiteral. got=%T", statement.Expression)
	}

	if len(list.Elements) != 3 {
		t.Fatalf("len(list.Elements) not 3. got=%d", len(list.Elements))
	}

	testNumberLiteral(t, list.Elements[0], 1)
	testInfixExpression(t, list.Elements[1], 2, "*", 2)
	testInfixExpression(t, list.Elements[2], 3, "+", 3)
}

func TestParsingMapLiteralStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	statement := program.Statements[0].(*ast.ExpressionStatement)
	mapLiteral, ok := statement.Expression.(*ast.MapLiteral)

	if !ok {
		t.Fatalf("expression is not ast.MapLiteral. got=%T", statement.Expression)
	}

	if len(mapLiteral.Pairs) != 3 {
		t.Errorf("mapLiteral.Pairs has wrong length. got=%d", len(mapLiteral.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range mapLiteral.Pairs {
		literal, ok := key.(*ast.StringLiteral)

		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}

		expectedValue := expected[literal.String()]

		testNumberLiteral(t, value, expectedValue)
	}
}

func TestParsingMapLiteralNumberKeys(t *testing.T) {
	input := `{1: 1, 2: 2, 3: 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	statement := program.Statements[0].(*ast.ExpressionStatement)
	mapLiteral, ok := statement.Expression.(*ast.MapLiteral)

	if !ok {
		t.Fatalf("expression is not ast.MapLiteral. got=%T", statement.Expression)
	}

	if len(mapLiteral.Pairs) != 3 {
		t.Errorf("mapLiteral.Pairs has wrong length. got=%d", len(mapLiteral.Pairs))
	}

	expected := map[string]int64{
		"1": 1,
		"2": 2,
		"3": 3,
	}

	for key, value := range mapLiteral.Pairs {
		literal, ok := key.(*ast.NumberLiteral)

		if !ok {
			t.Errorf("key is not ast.NumberLiteral. got=%T", key)
		}

		expectedValue := expected[literal.String()]

		testNumberLiteral(t, value, expectedValue)
	}
}

func TestParsingMapLiteralBooleanKeys(t *testing.T) {
	input := `{true: 1, false: 2}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	statement := program.Statements[0].(*ast.ExpressionStatement)
	mapLiteral, ok := statement.Expression.(*ast.MapLiteral)

	if !ok {
		t.Fatalf("expression is not ast.MapLiteral. got=%T", statement.Expression)
	}

	if len(mapLiteral.Pairs) != 2 {
		t.Errorf("mapLiteral.Pairs has wrong length. got=%d", len(mapLiteral.Pairs))
	}

	expected := map[string]int64{
		"true":  1,
		"false": 2,
	}

	for key, value := range mapLiteral.Pairs {
		boolean, ok := key.(*ast.BooleanLiteral)

		if !ok {
			t.Errorf("key is not ast.Boolean. got=%T", key)
		}

		expectedValue := expected[boolean.String()]

		testNumberLiteral(t, value, expectedValue)
	}
}

func TestParsingMapLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	statement := program.Statements[0].(*ast.ExpressionStatement)
	mapLiteral, ok := statement.Expression.(*ast.MapLiteral)

	if !ok {
		t.Fatalf("expression is not ast.MapLiteral. got=%T", statement.Expression)
	}

	if len(mapLiteral.Pairs) != 3 {
		t.Errorf("mapLiteral.Pairs has wrong length. got=%d", len(mapLiteral.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range mapLiteral.Pairs {
		literal, ok := key.(*ast.StringLiteral)

		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}

		testFunc, ok := tests[literal.String()]

		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}

		testFunc(value)
	}
}

func TestParsingEmptyMapLiteral(t *testing.T) {
	input := `{}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	statement := program.Statements[0].(*ast.ExpressionStatement)
	mapLiteral, ok := statement.Expression.(*ast.MapLiteral)

	if !ok {
		t.Fatalf("expression is not ast.MapLiteral. got=%T", statement.Expression)
	}

	if len(mapLiteral.Pairs) != 0 {
		t.Errorf("mapLiteral.Pairs has wrong length. got=%d", len(mapLiteral.Pairs))
	}
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myList[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	statement, _ := program.Statements[0].(*ast.ExpressionStatement)
	indexExpression, ok := statement.Expression.(*ast.IndexExpression)

	if !ok {
		t.Fatalf("expression not *ast.IndexExpression. got=%T", statement.Expression)
	}

	if !testIdentifierLiteral(t, indexExpression.Left, "myList") {
		return
	}

	if !testInfixExpression(t, indexExpression.Index, 1, "+", 1) {
		return
	}
}

func TestParsingProperties(t *testing.T) {
	input := "foo.bar"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	propertyExpression, ok := statement.Expression.(*ast.PropertyExpression)

	if !ok {
		t.Fatalf("expression not *ast.PropertyExpression. got=%T", statement.Expression)
	}

	if !testIdentifierLiteral(t, propertyExpression.Object, "foo") {
		return
	}

	if !testIdentifierLiteral(t, propertyExpression.Property, "bar") {
		return
	}
}

func TestWhileExpression(t *testing.T) {
	input := `while (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.WhileExpression)

	if !ok {
		t.Fatalf("statement.Expression is not ast.WhileExpression. got=%T", statement.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("while consequence is not 1 statement. got=%d", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Consequence.Statements[0] is not ast.ExpressionStatement. got=%T", expression.Consequence.Statements[0])
	}

	if !testIdentifierLiteral(t, consequence.Expression, "x") {
		return
	}
}

func testIdentifierLiteral(t *testing.T, expression ast.Expression, value string) bool {
	identifier, ok := expression.(*ast.IdentifierLiteral)
	if !ok {
		t.Errorf("expression not *ast.Identifier. got=%T", expression)
		return false
	}

	if identifier.Value != value {
		t.Errorf("identifier.Value not %s. got=%s", value, identifier.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, expression ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testNumberLiteral(t, expression, int64(v))
	case int64:
		return testNumberLiteral(t, expression, int64(v))
	case float64:
		return testNumberLiteral(t, expression, int64(v))
	case string:
		return testIdentifierLiteral(t, expression, v)
	case bool:
		return testBooleanLiteral(t, expression, v)
	}

	t.Errorf("type of expression not handled. got=%T", expression)

	return false
}

func testInfixExpression(t *testing.T, expression ast.Expression, left interface{}, operator string, right interface{}) bool {
	operatorExpression, ok := expression.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expression is not ast.InfixExpression. got=%T(%s)", expression, expression)
		return false
	}

	if !testLiteralExpression(t, operatorExpression.Left, left) {
		return false
	}

	if operatorExpression.Operator != operator {
		t.Errorf("expression.Operator is not '%s'. got=%q", operator, operatorExpression.Operator)
		return false
	}

	if !testLiteralExpression(t, operatorExpression.Right, right) {
		return false
	}

	return true
}

func TestForExpression(t *testing.T) {
	input := `for (x := 0; x < 10; x := x + 1) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.ForExpression)

	if !ok {
		t.Fatalf("statement.Expression is not ast.ForExpression. got=%T", statement.Expression)
	}

	if expression.Identifier != "x" {
		t.Errorf("wrong identifier in for loop. got=%s\n", expression.Identifier)
	}

	_, ok = expression.Initializer.(*ast.AssignStatement)

	if !ok {
		t.Fatalf("Initializer is not ast.AssignStatement. got=%T", expression.Initializer)
	}

	_, ok = expression.Increment.(*ast.AssignStatement)

	if !ok {
		t.Fatalf("Increment is not ast.AssignStatement. got=%T", expression.Increment)
	}

	_, ok = expression.Block.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement, got=%T", expression.Block.Statements[0])
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, message := range errors {
		t.Errorf("parser error: %q", message)
	}

	t.FailNow()
}
