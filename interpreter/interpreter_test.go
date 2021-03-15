package interpreter

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/environment"
	"ghostlang.org/x/ghost/glitch"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
	"github.com/shopspring/decimal"
)

func TestEvaluateNumberExpression(t *testing.T) {
	tests := []struct {
		literal  string
		expected interface{}
	}{
		{"5", 5},
		{"10", 10},
		{"3.14", 3.14},
		{"-5", -5},
		{"-10", -10},
		{"-3.14", -3.14},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, test := range tests {
		scanner := scanner.New(test.literal)
		tokens := scanner.ScanTokens()
		parser := parser.New(tokens)
		statements := parser.Parse()
		env := environment.New()

		if len(statements) != 1 {
			t.Fatalf("Expected 1 statement, got=%v", len(statements))
		}

		expression, ok := statements[0].(*ast.Expression)

		if !ok {
			t.Fatalf("Expected *ast.Expression, got=%T", statements[0])
		}

		result, _ := Evaluate(expression.Expression, env)

		verifyLiteralValue(result, test.expected, t)
	}
}

func TestEvaluateIfStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"if (true) { print 10 }", "10"},
		{"if (false) { print true } print false", "false"},
		{"if (1) { print 10 }", "10"},
		{"if (1 < 2) { print 10 }", "10"},
		{"if (1 > 2) { print true } print false", "false"},
		{"if (1 > 2) { print 10 } else { print 20 }", "20"},
		{"if (1 < 2) { print 10 } else { print 20 }", "10"},
		{"if (1 < 2) { print 10 } else if (1 == 1) { print 20 } else { print 30 }", "10"},
		{"if (1 > 2) {print  10 } else if (1 == 1) { print 20 } else { print 30 }", "20"},
		{"if (1 > 2) { print 10 } else if (1 == 2) { print 20 } else { print 30 }", "30"},
	}

	for _, test := range tests {
		result := new(bytes.Buffer)
		env := environment.New()
		env.SetWriter(result)

		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := parser.New(tokens)
		statements := parser.Parse()

		if glitch.HadParseError {
			return
		}

		Interpret(statements, env)

		if glitch.HadParseError || glitch.HadRuntimeError {
			os.Exit(1)
		}

		equals := strings.Compare(test.expected, string(bytes.TrimRight(result.Bytes(), "\n")))

		if equals != 0 {
			t.Errorf("expected value not %v, got=%v", test.expected, result.String())
		}
	}
}

func TestEvaluateWhileStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`
				let x = 5

				while (x < 10) {
					x = x + 1
				}

				print x
			`, "10",
		},
	}

	for _, test := range tests {
		result := new(bytes.Buffer)
		env := environment.New()
		env.SetWriter(result)

		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := parser.New(tokens)
		statements := parser.Parse()

		if glitch.HadParseError {
			return
		}

		Interpret(statements, env)

		if glitch.HadParseError || glitch.HadRuntimeError {
			os.Exit(1)
		}

		equals := strings.Compare(test.expected, string(bytes.TrimRight(result.Bytes(), "\n")))

		if equals != 0 {
			t.Errorf("expected value not %v, got=%v", test.expected, result.String())
		}
	}
}

// =============================================================================
// Helper methods

func verifyLiteralValue(literal interface{}, expected interface{}, t *testing.T) {
	switch result := literal.(type) {
	case *object.Number:
		verifyNumberValue(result, expected, t)
	// case bool:
	// 	verifyBooleanValue(result, expected, t)
	// case string:
	// 	verifyStringValue(result, expected, t)
	default:
		t.Fatalf("Unsupported literal type, expected float64, bool, or string, got=%T", result)
	}
}

func verifyNumberValue(number *object.Number, expected interface{}, t *testing.T) {
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

	equals := expected.(decimal.Decimal).Equal(number.Value)

	if !equals {
		t.Errorf("Numbers are not equal, expected %v, got=%v", expected, number.Value)
	}
}
