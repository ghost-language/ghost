package evaluator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/optimizer"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
	"ghostlang.org/x/ghost/token"
)

// Imports are process-wide: a module is read, evaluated, and kept, so that
// importing it from two places runs it once. That makes this shared state, and
// Ghost code can run on more than one goroutine at a time — an http.handle()
// callback runs per request — so it is guarded. An unguarded map written from
// two goroutines does not raise a Ghost error or even a Go panic; it kills the
// process outright, which is exactly the kind of failure this whole area exists
// to prevent.
var (
	moduleState sync.Mutex

	searchPaths []string

	// imported records every module that has been loaded. A module that is
	// still being evaluated is present with a nil scope, which is how a cycle
	// is recognised: a module cannot finish importing itself.
	imported = map[string]*object.Scope{}

	// loading names the modules currently being evaluated, innermost last, so a
	// circular import can describe the cycle rather than just refusing.
	loading []string
)

// rememberImport records a loaded module and reports what was already there.
func rememberImport(filename string, scope *object.Scope) {
	moduleState.Lock()
	defer moduleState.Unlock()

	imported[filename] = scope
}

// lookupImport reports a module that has already been loaded, and whether it
// has been seen at all — a module seen but still loading answers (nil, true).
func lookupImport(filename string) (*object.Scope, bool) {
	moduleState.Lock()
	defer moduleState.Unlock()

	scope, ok := imported[filename]

	return scope, ok
}

func evaluateImport(node *ast.Import, scope *object.Scope) object.Object {
	filename, err := resolveModule(node.Token, node.Path.Value, scope)

	if err != nil {
		return err
	}

	// Have we imported this file before? If so, we don't need to do anything.
	if _, ok := lookupImport(filename); ok {
		return nil
	}

	moduleScope, err := loadModule(filename, node.Token, node.Path.Value, scope)

	if err != nil {
		return err
	}

	rememberImport(filename, moduleScope)

	return nil
}

func evaluateImportFrom(node *ast.ImportFrom, scope *object.Scope) object.Object {
	filename, err := resolveModule(node.Token, node.Path.Value, scope)

	if err != nil {
		return err
	}

	moduleScope, loaded := lookupImport(filename)

	if !loaded {
		moduleScope, err = loadModule(filename, node.Token, node.Path.Value, scope)

		if err != nil {
			return err
		}

		rememberImport(filename, moduleScope)
	}

	if moduleScope == nil {
		return object.NewError(fault.Import, node.Token, "module `%s` is still being imported, so nothing can be taken from it yet", node.Path.Value).
			WithHelp("this is a circular import: %s", importCycle(node.Path.Value))
	}

	if node.Everything {
		for alias, value := range moduleScope.Environment.All() {
			scope.Environment.Set(alias, value)
		}

		return nil
	}

	for alias, identifier := range node.Identifiers {
		value, ok := moduleScope.Environment.Get(identifier.Value)

		if !ok {
			raised := object.NewError(fault.Import, node.Token, "module `%s` does not define `%s`", node.Path.Value, identifier.Value)

			if suggestion, ok := nearestName(identifier.Value, exported(moduleScope)); ok {
				raised.WithHelp("did you mean `%s`?", suggestion)
			}

			return raised
		}

		scope.Environment.Set(alias, value)
	}

	return nil
}

// resolveModule turns the name written in an import into the file it names.
func resolveModule(tok token.Token, name string, scope *object.Scope) (string, *object.Error) {
	addSearchPath(scope.Environment.GetDirectory())

	filename := findFile(name)

	if filename == "" {
		return "", object.NewError(fault.Import, tok, "no module found at `%s.ghost`", name).
			WithHelp("modules are looked for next to the file importing them")
	}

	return filename, nil
}

// loadModule reads, parses, and evaluates a module in a scope of its own.
//
// The module is marked as in progress before it runs, so a module that imports
// its way back to itself is caught here rather than looping until the stack
// runs out.
func loadModule(filename string, tok token.Token, name string, scope *object.Scope) (*object.Scope, *object.Error) {
	rememberImport(filename, nil)

	moduleState.Lock()
	loading = append(loading, name)
	moduleState.Unlock()

	defer func() {
		moduleState.Lock()
		defer moduleState.Unlock()

		loading = loading[:len(loading)-1]
	}()

	text, failure := os.ReadFile(filename)

	if failure != nil {
		return nil, object.NewError(fault.System, tok, "could not read module `%s`: %s", name, failure)
	}

	directory := scope.Environment.GetDirectory()
	currentFile := strings.Replace(filename, directory+"/", "", 1)

	parsed := parser.New(scanner.New(string(text), currentFile))
	program := parsed.Parse()

	// A module that will not parse cannot be imported. The first problem is
	// reported as the import's failure, carrying its own position inside the
	// module, with a frame naming the import that pulled it in.
	if raised := parsed.Errors(); len(raised) != 0 {
		return nil, object.NewErrorFrom(raised[0]).WithFrame(fmt.Sprintf("import of `%s`", name), tok)
	}

	moduleScope := &object.Scope{Self: scope.Self, Environment: object.NewEnvironment()}
	moduleScope.Environment.SetDirectory(directory)

	result := Evaluate(optimizer.Optimize(program), moduleScope)

	if failed, ok := result.(*object.Error); ok {
		return nil, failed.WithFrame(fmt.Sprintf("import of `%s`", name), tok)
	}

	return moduleScope, nil
}

// exported lists the names a module defines, for suggesting the one that was
// probably meant.
func exported(moduleScope *object.Scope) []string {
	names := make([]string, 0, 16)

	for name := range moduleScope.Environment.All() {
		names = append(names, name)
	}

	return names
}

func findFile(name string) string {
	moduleState.Lock()
	defer moduleState.Unlock()

	basename := fmt.Sprintf("%s.ghost", name)

	for _, path := range searchPaths {
		file := filepath.Join(path, basename)

		if fileExists(file) {
			return file
		}
	}

	return ""
}

func addSearchPath(path string) {
	moduleState.Lock()
	defer moduleState.Unlock()

	for _, existing := range searchPaths {
		if existing == path {
			return
		}
	}

	searchPaths = append(searchPaths, path)
}

// importCycle describes the chain of imports that led back to a module.
func importCycle(name string) string {
	moduleState.Lock()
	defer moduleState.Unlock()

	chain := make([]string, 0, len(loading)+1)
	chain = append(chain, loading...)
	chain = append(chain, name)

	return strings.Join(chain, " imports ")
}

func fileExists(file string) bool {
	info, err := os.Stat(file)

	return err == nil && !info.IsDir()
}
