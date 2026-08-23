package library

import (
	"testing"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

func testFunction(value string) object.GoFunction {
	return func(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
		return &object.String{Value: value}
	}
}

// An embedding host claims its own scheme rather than borrowing the standard
// library's "ghost" one, so its scripts get a namespace that reads as the
// host's own - Lumen's "lumen:font", say.
func TestRegisterModuleForSchemeIsIsolatedFromGhost(t *testing.T) {
	methods := map[string]*object.LibraryFunction{
		"name": {Name: "name", Function: testFunction("Lumen Sans")},
	}

	RegisterModuleForScheme("lumen-test", "font", methods, map[string]*object.LibraryProperty{})

	if _, ok := Modules["font"]; ok {
		t.Error("a scheme-registered module leaked into the standard library's own registry")
	}

	registry := Scheme("lumen-test")

	module, ok := registry.Modules["font"]

	if !ok {
		t.Fatal("module not found under its own scheme")
	}

	if module.Name != "font" {
		t.Errorf("got name=%q", module.Name)
	}
}

func TestRegisterFunctionForSchemeIsIsolatedFromGhost(t *testing.T) {
	RegisterFunctionForScheme("lumen-test", "greet", testFunction("hi"))

	if _, ok := Functions["greet"]; ok {
		t.Error("a scheme-registered function leaked into the standard library's own registry")
	}

	registry := Scheme("lumen-test")

	if _, ok := registry.Functions["greet"]; !ok {
		t.Fatal("function not found under its own scheme")
	}
}

func TestSchemesListsEveryRegisteredScheme(t *testing.T) {
	RegisterFunctionForScheme("lumen-schemes-test", "greet", testFunction("hi"))

	found := false

	for _, scheme := range Schemes() {
		if scheme == StandardScheme {
			found = true
		}
	}

	if !found {
		t.Error("StandardScheme (\"ghost\") should always be listed")
	}

	found = false

	for _, scheme := range Schemes() {
		if scheme == "lumen-schemes-test" {
			found = true
		}
	}

	if !found {
		t.Error("a scheme with a registration should be listed")
	}
}

// Scheme() must not make an unregistered scheme start looking registered:
// a lookup that finds nothing is how "wrong name" gets told apart from
// "wrong/not-yet-loaded scheme" one level up, in evaluator/import.go.
func TestSchemeOnAnUnregisteredNameDoesNotAppearInSchemes(t *testing.T) {
	registry := Scheme("nobody-has-registered-this")

	if len(registry.Modules) != 0 || len(registry.Functions) != 0 {
		t.Fatal("expected an empty registry for an unregistered scheme")
	}

	for _, scheme := range Schemes() {
		if scheme == "nobody-has-registered-this" {
			t.Error("an unregistered scheme should not appear in Schemes()")
		}
	}
}

func TestLocateFindsEveryOwningScheme(t *testing.T) {
	RegisterFunctionForScheme("lumen-locate-a", "widget", testFunction("a"))
	RegisterFunctionForScheme("lumen-locate-b", "widget", testFunction("b"))

	schemes := Locate("widget")

	if len(schemes) != 2 {
		t.Fatalf("expected 2 owning schemes, got %v", schemes)
	}

	if schemes[0] != "lumen-locate-a" || schemes[1] != "lumen-locate-b" {
		t.Errorf("expected sorted [lumen-locate-a lumen-locate-b], got %v", schemes)
	}
}

func TestLocateFindsStandardLibraryModules(t *testing.T) {
	schemes := Locate("math")

	if len(schemes) != 1 || schemes[0] != StandardScheme {
		t.Errorf("expected [%s], got %v", StandardScheme, schemes)
	}
}

func TestLocateReportsNothingForAnUnregisteredName(t *testing.T) {
	if schemes := Locate("thisNameIsNotRegisteredAnywhere"); len(schemes) != 0 {
		t.Errorf("expected no owning schemes, got %v", schemes)
	}
}

// RegisterModule/RegisterFunction — the pre-existing, unscoped embedding API
// — keep registering into the standard library's own "ghost" scheme, so
// existing embedding code is unaffected by scheme registration existing.
func TestUnscopedRegisterFunctionsStillTargetStandardScheme(t *testing.T) {
	RegisterFunction("libraryTestBackwardCompat", testFunction("x"))

	if _, ok := Functions["libraryTestBackwardCompat"]; !ok {
		t.Fatal("RegisterFunction should still write into the package-level Functions map")
	}

	schemes := Locate("libraryTestBackwardCompat")

	if len(schemes) != 1 || schemes[0] != StandardScheme {
		t.Errorf("expected [%s], got %v", StandardScheme, schemes)
	}
}
