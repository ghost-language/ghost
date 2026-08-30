package optimizer

import (
	"testing"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
)

func optimizeSource(t *testing.T, source string) ast.Node {
	t.Helper()

	p := parser.New(scanner.New(source, "test.gs"))
	program := Optimize(p.Parse())

	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors for %q: %v", source, p.Errors())
	}

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement for %q, got %d", source, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.Expression)

	if !ok {
		t.Fatalf("statement is not ast.Expression for %q. got=%T", source, program.Statements[0])
	}

	return statement.Expression
}

func TestFoldsIntegerArithmetic(t *testing.T) {
	tests := []struct {
		source   string
		expected int64
	}{
		{"1 + 2", 3},
		{"10 - 4", 6},
		{"6 * 7", 42},
		{"7 % 3", 1},
		{"-5", -5},
		{"2 * 3 + 4 * 5 - 6", 20},
		{"1 + 2 + 3 + 4 + 5", 15},
		{"-(3 * 4)", -12},
	}

	for _, tt := range tests {
		folded := optimizeSource(t, tt.source)

		number, ok := folded.(*ast.Number)

		if !ok {
			t.Errorf("%q did not fold to a number. got=%T", tt.source, folded)
			continue
		}

		if number.IsFloat {
			t.Errorf("%q folded to a float, expected an integer", tt.source)
			continue
		}

		if number.IntValue != tt.expected {
			t.Errorf("%q folded to %d, expected %d", tt.source, number.IntValue, tt.expected)
		}
	}
}

func TestFoldsFloatArithmetic(t *testing.T) {
	tests := []struct {
		source   string
		expected float64
	}{
		{"1.5 + 2.5", 4},
		{"3.0 * 2.0", 6},
		{"1 + 2.5", 3.5},
		// Division always promotes to a float, matching object.Number.Div.
		{"6 / 3", 2},
		{"7 / 2", 3.5},
	}

	for _, tt := range tests {
		folded := optimizeSource(t, tt.source)

		number, ok := folded.(*ast.Number)

		if !ok {
			t.Errorf("%q did not fold to a number. got=%T", tt.source, folded)
			continue
		}

		if !number.IsFloat {
			t.Errorf("%q folded to an integer, expected a float", tt.source)
			continue
		}

		if number.FloatValue != tt.expected {
			t.Errorf("%q folded to %v, expected %v", tt.source, number.FloatValue, tt.expected)
		}
	}
}

func TestFoldsComparisonsAndBooleans(t *testing.T) {
	tests := []struct {
		source   string
		expected bool
	}{
		{"1 < 2", true},
		{"2 <= 2", true},
		{"3 > 4", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{`"a" == "a"`, true},
		{`"a" < "b"`, true},
		{"true and false", false},
		{"true or false", true},
		{"!true", false},
		{"!false", true},
		{"!null", true},
		{"!(1 == 2)", true},
		{"!5", false},

		// A string literal's folded truthiness depends on whether it's
		// empty (§8.5, §13.11) - this used to hard-code every string to
		// fold to false, regardless of content, silently diverging from
		// the evaluator's own isTruthy rule for the one case (an empty
		// string literal) where it mattered.
		{`!"abc"`, false},
		{`!""`, true},
	}

	for _, tt := range tests {
		folded := optimizeSource(t, tt.source)

		boolean, ok := folded.(*ast.Boolean)

		if !ok {
			t.Errorf("%q did not fold to a boolean. got=%T", tt.source, folded)
			continue
		}

		if boolean.Value != tt.expected {
			t.Errorf("%q folded to %t, expected %t", tt.source, boolean.Value, tt.expected)
		}
	}
}

func TestFoldsStringConcatenation(t *testing.T) {
	folded := optimizeSource(t, `"hello" + " " + "world"`)

	str, ok := folded.(*ast.String)

	if !ok {
		t.Fatalf("string concatenation did not fold. got=%T", folded)
	}

	if str.Value != "hello world" {
		t.Errorf("folded to %q, expected %q", str.Value, "hello world")
	}
}

// Expressions that would raise a runtime error, or that depend on values only
// known at run time, must be left alone so that the evaluator still sees them.
func TestLeavesUnfoldableExpressions(t *testing.T) {
	tests := []string{
		"1 / 0",        // division by zero is a runtime error
		"1 % 0",        // modulo by zero is a runtime error
		"1.0 / 0.0",    // float division by zero is a runtime error too
		"1 + true",     // type mismatch is a runtime error
		`1 + "a"`,      // type mismatch is a runtime error
		"1 .. 5",       // a range builds a mutable list object
		"a + 1",        // depends on a runtime variable
		"-true",        // negating a non-number is a runtime error
		"foo() + 1",    // calls may have side effects
		"[1, 2] + [3]", // lists are not folded
	}

	for _, source := range tests {
		folded := optimizeSource(t, source)

		switch folded.(type) {
		case *ast.Number, *ast.String, *ast.Boolean:
			t.Errorf("%q was folded to a literal (%T) but should have been left alone", source, folded)
		}
	}
}

// Folding must reach into every construct that contains expressions, not just
// the top level of a program.
func TestFoldsInsideNestedConstructs(t *testing.T) {
	source := `function outer() {
		if (1 + 1 == 2) {
			for (i = 0; i < 2 * 5; i = i + 1) {
				total = total + (3 * 4)
			}
		}
		return 8 - 3
	}`

	p := parser.New(scanner.New(source, "test.gs"))
	program := Optimize(p.Parse())

	if len(p.Errors()) != 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	// Walk the tree and assert no foldable infix expression survives.
	var unfolded []string

	var walk func(node ast.Node)
	walk = func(node ast.Node) {
		switch node := node.(type) {
		case *ast.Program:
			for _, s := range node.Statements {
				walk(s)
			}
		case *ast.Block:
			for _, s := range node.Statements {
				walk(s)
			}
		case *ast.Expression:
			walk(node.Expression)
		case *ast.Function:
			walk(node.Body)
		case *ast.If:
			walk(node.Condition)
			walk(node.Consequence)
		case *ast.For:
			walk(node.Condition)
			walk(node.Block)
		case *ast.Assign:
			walk(node.Value)
		case *ast.Return:
			walk(node.Value)
		case *ast.Infix:
			if foldInfix(node) != nil {
				unfolded = append(unfolded, "infix survived folding")
			}
			walk(node.Left)
			walk(node.Right)
		}
	}

	walk(program)

	if len(unfolded) != 0 {
		t.Errorf("found %d foldable expressions that were not folded", len(unfolded))
	}
}
