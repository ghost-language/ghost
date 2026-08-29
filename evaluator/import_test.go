package evaluator

import (
	"os"
	"path/filepath"
	"testing"

	"ghostlang.org/x/ghost/library/modules"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
)

// evaluateInDirectory runs input as if it were a file living in dir, so an
// `import` inside it resolves modules written alongside it.
func evaluateInDirectory(dir string, input string) object.Object {
	scope := &object.Scope{Environment: object.NewEnvironment()}
	scope.Environment.SetDirectory(dir)

	object.RegisterEvaluator(Evaluate)
	modules.RegisterEvaluator(Evaluate)

	p := parser.New(scanner.New(input, "main.gs"))
	program := p.Parse()

	return Evaluate(program, scope)
}

func writeModule(t *testing.T, dir, name, source string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, name), []byte(source), 0o644); err != nil {
		t.Fatalf("failed to write module %s: %v", name, err)
	}
}

func TestImportBindsDefaultName(t *testing.T) {
	dir := t.TempDir()

	writeModule(t, dir, "greeting.gs", `hello = "hi"`)

	result := evaluateInDirectory(dir, "import \"greeting\"\ngreeting.hello")

	isStringObject(t, result, "hi")
}

func TestImportAsBindsAlias(t *testing.T) {
	dir := t.TempDir()

	writeModule(t, dir, "greeting.gs", `hello = "hi"`)

	result := evaluateInDirectory(dir, "import \"greeting\" as g\ng.hello")

	isStringObject(t, result, "hi")
}

func TestImportFromWithBraces(t *testing.T) {
	dir := t.TempDir()

	writeModule(t, dir, "greeting.gs", "hello = \"hi\"\nbye = \"bye\"")

	result := evaluateInDirectory(dir, "import { hello, bye } from \"greeting\"\nhello + \" \" + bye")

	isStringObject(t, result, "hi bye")
}

func TestImportFromWithBracesAndAlias(t *testing.T) {
	dir := t.TempDir()

	writeModule(t, dir, "greeting.gs", `hello = "hi"`)

	result := evaluateInDirectory(dir, "import { hello as h } from \"greeting\"\nh")

	isStringObject(t, result, "hi")
}

func TestImportCombinedModuleAndNamed(t *testing.T) {
	dir := t.TempDir()

	writeModule(t, dir, "greeting.gs", `hello = "hi"`)

	result := evaluateInDirectory(dir, "import greeting, { hello } from \"greeting\"\ngreeting.hello + \" \" + hello")

	isStringObject(t, result, "hi hi")
}

func TestImportCombinedSchemeModuleAndNamed(t *testing.T) {
	dir := t.TempDir()

	result := evaluateInDirectory(dir, "import math, { pi } from \"ghost:math\"\nmath.pi == pi")

	isBooleanObject(t, result, true)
}

func TestImportReimportStillBinds(t *testing.T) {
	dir := t.TempDir()

	writeModule(t, dir, "shared.gs", `value = 1`)

	// Importing the same module twice from two different places must not skip
	// the second binding just because the module itself only needs to run
	// once.
	result := evaluateInDirectory(dir, "import \"shared\"\nimport \"shared\"\nshared.value")

	isNumberObject(t, result, 1)
}
