package interpreter

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"ghostlang.org/x/ghost/errors"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
)

func TestNamedFunctionCalls(t *testing.T) {
	tests := []struct {
		input string
		expected string
	}{
		{"function five() { print 5 } five()", "5"},
		{"function ten() { print 10 } ten()", "10"},
		{"function foo() { print \"foo\" } foo()", "foo"},
	}

	for _, test := range tests {
		result := new(bytes.Buffer)
		env := object.NewEnvironment()
		env.SetWriter(result)

		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := parser.New(tokens)
		statements := parser.Parse()

		if errors.HadParseError {
			return
		}

		Interpret(statements, env)

		if errors.HadParseError || errors.HadRuntimeError {
			os.Exit(1)
		}

		equals := strings.Compare(test.expected, string(bytes.TrimRight(result.Bytes(), "\n")))

		if equals != 0 {
			t.Errorf("expected value not %v, got=%v", test.expected, result.String())
		}
	}
}