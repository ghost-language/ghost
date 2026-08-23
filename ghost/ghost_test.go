package ghost

import (
	"bytes"
	"strings"
	"testing"

	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

func TestExecuteReportsSyntaxErrorsInFull(t *testing.T) {
	report := &bytes.Buffer{}

	instance := New()
	instance.SetReportWriter(report)
	instance.SetFile("test.ghost")
	instance.SetSource("x = 1\ny = )\n")

	result := instance.Execute()

	if !object.IsError(result) {
		t.Fatalf("expected an error, got %T", result)
	}

	written := report.String()

	for _, expected := range []string{"syntax error", "test.ghost:2:5", "y = )", "^"} {
		if !strings.Contains(written, expected) {
			t.Errorf("report is missing %q:\n%s", expected, written)
		}
	}
}

func TestExecuteReportsRuntimeErrorsInFull(t *testing.T) {
	report := &bytes.Buffer{}

	instance := New()
	instance.SetReportWriter(report)
	instance.SetFile("test.ghost")
	instance.SetSource("count = 3\ntotal = count + \"items\"\n")

	instance.Execute()

	written := report.String()

	for _, expected := range []string{"type error", "test.ghost:2:15", "total = count + \"items\""} {
		if !strings.Contains(written, expected) {
			t.Errorf("report is missing %q:\n%s", expected, written)
		}
	}
}

// An embedder usually wants to present failures itself. It can: the error comes
// back either way, carrying everything a report is built from.
func TestQuietStillReturnsTheError(t *testing.T) {
	report := &bytes.Buffer{}

	instance := New()
	instance.SetReportWriter(report)
	instance.SetQuiet(true)
	instance.SetFile("test.ghost")
	instance.SetSource("missing")

	result := instance.Execute()

	if report.Len() != 0 {
		t.Errorf("expected nothing to be written, got:\n%s", report)
	}

	raised, ok := result.(*object.Error)

	if !ok {
		t.Fatalf("expected an error, got %T", result)
	}

	if raised.Fault.Kind != fault.Name {
		t.Errorf("got kind=%v, expected name", raised.Fault.Kind)
	}
}

// Nothing beneath Execute is meant to panic, but "meant to" is not "does" — and
// an embedder's own Go code is registered right alongside Ghost's. Either way a
// panic must reach the reader as a Ghost error rather than as a Go traceback
// naming files they have never heard of.
func TestAPanicBecomesAnInternalError(t *testing.T) {
	RegisterFunction("explode", func(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
		panic("boom")
	})

	instance := New()
	instance.SetQuiet(true)
	instance.SetFile("test.ghost")
	instance.SetSource("import \"ghost:explode\"\nexplode()")

	result := instance.Execute()

	raised, ok := result.(*object.Error)

	if !ok {
		t.Fatalf("expected an error, got %T", result)
	}

	if raised.Fault.Kind != fault.Internal {
		t.Errorf("got kind=%v, expected internal", raised.Fault.Kind)
	}

	if !strings.Contains(raised.Message(), "boom") {
		t.Errorf("got=%q", raised.Message())
	}

	if raised.Fault.Help == "" {
		t.Error("an internal error should ask to be reported")
	}
}

// An embedding host isn't limited to the standard library's own `ghost:`
// import scheme — it can claim one of its own, so its scripts read as its
// own rather than borrowing Ghost's namespace (a Go program embedding Ghost
// as "Lumen" registering "font" under "lumen:", say).
func TestEmbedderCanRegisterAModuleUnderItsOwnScheme(t *testing.T) {
	methods := map[string]*object.LibraryFunction{
		"name": {Name: "name", Function: func(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
			return &object.String{Value: "Lumen Sans"}
		}},
	}

	RegisterModuleForScheme("lumen", "font", methods, map[string]*object.LibraryProperty{})

	instance := New()
	instance.SetQuiet(true)
	instance.SetFile("test.ghost")
	instance.SetSource("import \"lumen:font\"\nfont.name()")

	result := instance.Execute()

	str, ok := result.(*object.String)

	if !ok {
		t.Fatalf("expected a string, got %T (%v)", result, result)
	}

	if str.Value != "Lumen Sans" {
		t.Errorf("got=%q", str.Value)
	}
}

// A standalone function registered under a custom scheme is reached the
// bare way (`import "lumen:greet"`), the same as a standalone standard
// library function (`import "ghost:type"` would work identically) — not
// through the `from` form, which is for pulling members out of a module.
func TestEmbedderCanRegisterAFunctionUnderItsOwnScheme(t *testing.T) {
	RegisterFunctionForScheme("lumen", "greet", func(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
		return &object.String{Value: "hi from lumen"}
	})

	instance := New()
	instance.SetQuiet(true)
	instance.SetFile("test.ghost")
	instance.SetSource("import \"lumen:greet\"\ngreet()")

	result := instance.Execute()

	str, ok := result.(*object.String)

	if !ok {
		t.Fatalf("expected a string, got %T (%v)", result, result)
	}

	if str.Value != "hi from lumen" {
		t.Errorf("got=%q", str.Value)
	}
}

// Importing from a scheme nothing has ever registered under is a distinct,
// named failure - not a generic "not found" indistinguishable from a
// misspelled name within a real scheme.
func TestImportingAnUnregisteredSchemeIsReported(t *testing.T) {
	instance := New()
	instance.SetQuiet(true)
	instance.SetFile("test.ghost")
	instance.SetSource("import \"nosuchscheme:thing\"")

	result := instance.Execute()

	raised, ok := result.(*object.Error)

	if !ok {
		t.Fatalf("expected an error, got %T", result)
	}

	if raised.Fault.Kind != fault.Import {
		t.Errorf("got kind=%v, expected import", raised.Fault.Kind)
	}
}
