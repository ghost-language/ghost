package ghost

import (
	"bytes"
	"strings"
	"testing"

	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

// TestExecuteReusesScopeAcrossCalls confirms the mechanism `-i` (§13.4,
// cmd/ghost.go) relies on: a second Execute() on the same instance sees
// everything the first one bound, rather than starting from a fresh scope.
// This is what lets the CLI hand an already-run script's instance straight
// to the REPL and have its variables still be there to inspect.
func TestExecuteReusesScopeAcrossCalls(t *testing.T) {
	instance := New()
	instance.SetQuiet(true)
	instance.SetFile("test.gs")

	instance.SetSource("x = 41")

	if result := instance.Execute(); object.IsError(result) {
		t.Fatalf("unexpected error on first Execute: %s", result.String())
	}

	instance.SetSource("x + 1")

	result := instance.Execute()

	if object.IsError(result) {
		t.Fatalf("unexpected error on second Execute: %s", result.String())
	}

	number, ok := result.(*object.Number)

	if !ok {
		t.Fatalf("result is not Number. got=%T (%+v)", result, result)
	}

	if number.Int64() != 42 {
		t.Errorf("x + 1: got=%d, expected=42", number.Int64())
	}
}

func TestExecuteReportsSyntaxErrorsInFull(t *testing.T) {
	report := &bytes.Buffer{}

	instance := New()
	instance.SetReportWriter(report)
	instance.SetFile("test.gs")
	instance.SetSource("x = 1\ny = )\n")

	result := instance.Execute()

	if !object.IsError(result) {
		t.Fatalf("expected an error, got %T", result)
	}

	written := report.String()

	for _, expected := range []string{"syntax error", "test.gs:2:5", "y = )", "^"} {
		if !strings.Contains(written, expected) {
			t.Errorf("report is missing %q:\n%s", expected, written)
		}
	}
}

func TestExecuteReportsRuntimeErrorsInFull(t *testing.T) {
	report := &bytes.Buffer{}

	instance := New()
	instance.SetReportWriter(report)
	instance.SetFile("test.gs")
	instance.SetSource("count = 3\ntotal = count + \"items\"\n")

	instance.Execute()

	written := report.String()

	for _, expected := range []string{"type error", "test.gs:2:15", "total = count + \"items\""} {
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
	instance.SetFile("test.gs")
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
	instance.SetFile("test.gs")
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
	instance.SetFile("test.gs")
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
	instance.SetFile("test.gs")
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

// audioInstance stands in for what a real embedder's native instance looks
// like: its own Go type, its own state, its own Method() dispatch - nothing
// about it is Ghost-specific beyond implementing object.Object, exactly as
// object/native_class.go promises a NativeClass's Constructor is free to
// return.
type audioInstance struct {
	path string
}

func (audio *audioInstance) String() string    { return "audio instance" }
func (audio *audioInstance) Type() object.Type { return object.INSTANCE }

func (audio *audioInstance) Method(method string, tok token.Token, args []object.Object) (object.Object, bool) {
	if method == "path" {
		return &object.String{Value: audio.path}, true
	}

	return nil, false
}

func audioConstructor(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := object.Arity("Audio()", tok, args, 1); err != nil {
		return err
	}

	path, err := object.StringArgument("Audio()", tok, args, 0)

	if err != nil {
		return err
	}

	return &audioInstance{path: path.Value}
}

// A native class - the Lumen "Audio" example from the design discussion -
// is `new`-able and its instances behave like any other object, once
// registered as a member of a module the same way a method or property is.
func TestNativeClassCanBeNewedAfterNamedImport(t *testing.T) {
	RegisterClassForScheme("lumen", "audio", "Audio", audioConstructor)

	instance := New()
	instance.SetQuiet(true)
	instance.SetFile("test.gs")
	instance.SetSource(`import { Audio } from "lumen:audio"` + "\n" + `new Audio("path/to/file.mp3").path()`)

	result := instance.Execute()

	str, ok := result.(*object.String)

	if !ok {
		t.Fatalf("expected a string, got %T (%v)", result, result)
	}

	if str.Value != "path/to/file.mp3" {
		t.Errorf("got=%q", str.Value)
	}
}

// A class reached through the whole-module bare import, then dotted access,
// resolves identically - `audio.Audio` reads a class off a module exactly
// like it reads a method or a property.
func TestNativeClassCanBeNewedViaModuleDotAccess(t *testing.T) {
	RegisterClassForScheme("lumen", "audio", "Audio", audioConstructor)

	instance := New()
	instance.SetQuiet(true)
	instance.SetFile("test.gs")
	instance.SetSource(`import "lumen:audio"` + "\n" + `new audio.Audio("x.mp3").path()`)

	result := instance.Execute()

	str, ok := result.(*object.String)

	if !ok {
		t.Fatalf("expected a string, got %T (%v)", result, result)
	}

	if str.Value != "x.mp3" {
		t.Errorf("got=%q", str.Value)
	}
}

// Calling or reading off a native class itself (not an instance) is refused
// with the same help a Ghost-defined class already gives - a script author
// hitting this shouldn't be able to tell which kind of class it is.
func TestNativeClassRejectsMethodAndPropertyOnTheClassItself(t *testing.T) {
	RegisterClassForScheme("lumen", "audio", "Audio", audioConstructor)

	tests := []struct {
		name     string
		source   string
		expected string
	}{
		{"method call on the class", `import { Audio } from "lumen:audio"` + "\n" + `Audio.path()`, "class `Audio` has no method `path` to call on the class itself"},
		{"property read on the class", `import { Audio } from "lumen:audio"` + "\n" + `Audio.path`, "class `Audio` has no property `path` to read on the class itself"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance := New()
			instance.SetQuiet(true)
			instance.SetFile("test.gs")
			instance.SetSource(tt.source)

			result := instance.Execute()

			raised, ok := result.(*object.Error)

			if !ok {
				t.Fatalf("expected an error, got %T", result)
			}

			if !strings.Contains(raised.Message(), tt.expected) {
				t.Errorf("got=%q, expected to contain %q", raised.Message(), tt.expected)
			}
		})
	}
}

// A misspelled named import of a class gets the same nearest-name
// suggestion a misspelled method or property does.
func TestNativeClassTypoSuggestsTheClassName(t *testing.T) {
	RegisterClassForScheme("lumen", "audio", "Audio", audioConstructor)

	instance := New()
	instance.SetQuiet(true)
	instance.SetFile("test.gs")
	instance.SetSource(`import { Audoi } from "lumen:audio"`)

	result := instance.Execute()

	raised, ok := result.(*object.Error)

	if !ok {
		t.Fatalf("expected an error, got %T", result)
	}

	if raised.Fault.Help != "did you mean `Audio`?" {
		t.Errorf("got help=%q", raised.Fault.Help)
	}
}

// Importing from a scheme nothing has ever registered under is a distinct,
// named failure - not a generic "not found" indistinguishable from a
// misspelled name within a real scheme.
func TestImportingAnUnregisteredSchemeIsReported(t *testing.T) {
	instance := New()
	instance.SetQuiet(true)
	instance.SetFile("test.gs")
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
