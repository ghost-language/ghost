package compiler

import (
	"fmt"
	"testing"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/code"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}

func TestNumberArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
			},
		},
	}

	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		compiler := New()

		err := compiler.Compile(program)

		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		bytecode := compiler.Bytecode()

		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)

		if err != nil {
			t.Errorf("testInstructions failed: %s", err)
		}

		err = testConstants(tt.expectedConstants, bytecode.Constants)

		if err != nil {
			t.Errorf("testConstants failed: %s", err)
		}
	}
}

func parse(input string) *ast.Program {
	scanner := scanner.New(input, "test.ghost")
	parser := parser.New(scanner)

	return parser.Parse()
}

func testInstructions(expected []code.Instructions, actual code.Instructions) error {
	concatted := concatInstructions(expected)

	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length.\nwant=%q\ngot=%q", concatted, actual)
	}

	for i, instruction := range concatted {
		if actual[i] != instruction {
			return fmt.Errorf("wrong instruction at %d.\nwant=%q\ngot=%q", i, concatted, actual)
		}
	}

	return nil
}

func testConstants(expected []interface{}, actual []object.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got=%d, want=%d", len(actual), len(expected))
	}

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testNumberObject(int64(constant), actual[i])

			if err != nil {
				return fmt.Errorf("constant %d - testNumberObject failed: %s", i, err)
			}
		}
	}

	return nil
}

func concatInstructions(instructions []code.Instructions) code.Instructions {
	out := code.Instructions{}

	for _, instruction := range instructions {
		out = append(out, instruction...)
	}

	return out
}

func testNumberObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Number)

	if !ok {
		return fmt.Errorf("object is not Number. got=%T (%+v)", actual, actual)
	}

	if result.Value.IntPart() != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}

	return nil
}
