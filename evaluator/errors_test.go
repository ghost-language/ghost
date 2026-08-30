package evaluator

import (
	"strings"
	"testing"

	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
)

// A name that is not defined is nearly always a name that was mistyped, and the
// interpreter is holding the list of names that were in scope when it gave up.
func TestUndefinedNamesSuggestTheNearest(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected string
	}{
		{"a transposed variable", `name = "ghost" type(nmae)`, "did you mean `name`?"},
		{"a dropped letter", `total = 1 type(totl)`, "did you mean `total`?"},
		{"a library function", `typ(5)`, "did you mean `type`?"},
		{"a library module", `mathh.abs(-1)`, "did you mean `math`? import it: `import \"ghost:math\"`"},
		{"a method on a class", `class Point { distance() { return 1 } } new Point().distence()`, "did you mean `distance`?"},
		{"a module method", "import \"ghost:math\"\nmath.sqrtt(4)", "did you mean `math.sqrt()`?"},
		{"a module property", "import \"ghost:math\"\nmath.pii", "did you mean `math.pi`?"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			raised := errorFrom(t, test.source)

			if raised.Fault.Help != test.expected {
				t.Errorf("got=%q, expected=%q", raised.Fault.Help, test.expected)
			}
		})
	}
}

// Guessing wrongly is worse than not guessing: a name with nothing like it in
// scope gets no suggestion at all.
func TestUnrelatedNamesGetNoSuggestion(t *testing.T) {
	raised := errorFrom(t, `name = "ghost" type(quixotic)`)

	if raised.Fault.Help != "" {
		t.Errorf("expected no suggestion, got %q", raised.Fault.Help)
	}
}

// Ghost runs on the Go stack, and a Go stack overflow cannot be recovered from:
// it kills the process outright. Counting the depth turns a runaway recursion
// back into an ordinary error.
func TestRunawayRecursionIsReported(t *testing.T) {
	raised := errorFrom(t, `function loop(n) { return loop(n + 1) } loop(0)`)

	if raised.Fault.Kind != fault.Value {
		t.Errorf("got kind=%v, expected value", raised.Fault.Kind)
	}

	if !strings.Contains(raised.Message(), "call depth") {
		t.Errorf("got=%q", raised.Message())
	}
}

// A failure inside a helper is much easier to place when the report says how
// the program got there.
func TestErrorsCarryTheCallsTheyPassedThrough(t *testing.T) {
	raised := errorFrom(t, "function inner() { return missing }\nfunction outer() { return inner() }\nouter()")

	if len(raised.Fault.Trace) != 2 {
		t.Fatalf("got %d frames, expected 2: %v", len(raised.Fault.Trace), raised.Fault.Trace)
	}

	if raised.Fault.Trace[0].Name != "inner()" || raised.Fault.Trace[1].Name != "outer()" {
		t.Errorf("frames are in the wrong order: %v", raised.Fault.Trace)
	}
}

// A library function reports its own failures at the call itself, so a frame
// there would only repeat the position already printed above it.
func TestLibraryCallsDoNotAddAFrame(t *testing.T) {
	raised := errorFrom(t, `type(1, 2)`)

	if len(raised.Fault.Trace) != 0 {
		t.Errorf("expected no frames, got %v", raised.Fault.Trace)
	}
}

// Methods on values used to assert their arguments' types inline, so calling
// one wrongly took the whole interpreter down with a Go panic.
func TestValueMethodsReportBadArguments(t *testing.T) {
	tests := []struct {
		source   string
		expected string
	}{
		{`"a,b".split()`, "test.gs:1:7: argument error: `string.split()` expects 1 argument, got 0"},
		{`"a,b".split(5)`, "test.gs:1:7: argument error: `string.split()` expects argument 1 to be a string, got number"},
		{`[1, 2].join()`, "test.gs:1:8: argument error: `list.join()` expects 1 argument, got 0"},
		{`[1, 2].join(5)`, "test.gs:1:8: argument error: `list.join()` expects argument 1 to be a string, got number"},
		{`[1, 2].push()`, "test.gs:1:8: argument error: `list.push()` expects 1 argument, got 0"},
		{`[1, 2].concat(3)`, "test.gs:1:8: argument error: `list.concat()` expects argument 1 to be a list, got number"},
		{`[1, 2].length(3)`, "test.gs:1:8: argument error: `list.length()` expects 0 arguments, got 1"},
		{`(1.5).round("two")`, "test.gs:1:7: argument error: `number.round()` expects argument 1 to be a number, got string"},
	}

	for _, test := range tests {
		raised := errorFrom(t, test.source)

		if raised.String() != test.expected {
			t.Errorf("got=%q, expected=%q", raised.String(), test.expected)
		}
	}
}

// A string is only a pattern when a method treats it as one, so a pattern that
// will not compile is the caller's mistake rather than a crash.
func TestBadPatternsAreReported(t *testing.T) {
	raised := errorFrom(t, `"x".matches("(")`)

	if raised.Fault.Kind != fault.Value {
		t.Errorf("got kind=%v, expected value", raised.Fault.Kind)
	}

	if !strings.Contains(raised.Message(), "as a pattern") {
		t.Errorf("got=%q", raised.Message())
	}
}

// Every error carries the kind of mistake it is, so a reader can tell a typo
// from a bad argument from a missing file before reading the sentence.
func TestErrorsAreClassified(t *testing.T) {
	tests := []struct {
		source   string
		expected fault.Kind
	}{
		{`missing`, fault.Name},
		{`1 + "a"`, fault.Type},
		{`type(1, 2)`, fault.Argument},
		{`1 / 0`, fault.Value},
		{`5.nope()`, fault.Property},
		{`import "not-a-real-module"`, fault.Import},
		{"import \"ghost:file\"\nfile.read(\"not-a-real-file\")", fault.System},
	}

	for _, test := range tests {
		raised := errorFrom(t, test.source)

		if raised.Fault.Kind != test.expected {
			t.Errorf("%q: got kind=%v, expected=%v", test.source, raised.Fault.Kind, test.expected)
		}
	}
}

// errorFrom evaluates a program that is expected to fail and returns the error.
func errorFrom(t *testing.T, source string) *object.Error {
	t.Helper()

	result := evaluate(source)

	raised, ok := result.(*object.Error)

	if !ok {
		t.Fatalf("%q: expected an error, got %T (%v)", source, result, result)
	}

	return raised
}
