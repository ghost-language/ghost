package evaluator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/library"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/optimizer"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
	"ghostlang.org/x/ghost/token"
)

// schemePattern recognizes a `scheme:` prefix on an import path — `ghost:`
// for the standard library, or one a Go host embedding Ghost claims for
// itself (Lumen's `lumen:`, say — see library.RegisterModuleForScheme). Two
// or more letters are required before the `:` specifically so a Windows
// drive letter (`C:\...`) is never mistaken for one; nothing in Ghost's own
// import paths uses a single-letter prefix.
var schemePattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+.-]+:`)

// schemeImport splits an import path into its scheme and the name after it,
// if it has one: `"ghost:math"` splits to `("ghost", "math")`. A path with
// no such prefix is a `.ghost` file import instead (§8.9) — the two are told
// apart by this alone, no registration lookup needed to decide which.
func schemeImport(path string) (scheme string, name string, ok bool) {
	prefix := schemePattern.FindString(path)

	if prefix == "" {
		return "", "", false
	}

	return prefix[:len(prefix)-1], path[len(prefix):], true
}

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

// evaluateImport handles the bare form of an import: `import "math"` binds
// the whole module to a variable named after its path, and `import "math" as
// m` binds it to `m` instead. The module is only read, parsed, and evaluated
// once no matter how many places import it (loadModule/rememberImport below
// are shared with evaluateImportFrom for that reason); every import of it
// still has to bind its own name into the importing scope.
func evaluateImport(node *ast.Import, scope *object.Scope) object.Object {
	if scheme, name, ok := schemeImport(node.Path.Value); ok {
		return evaluateSchemeImport(scheme, name, node, scope)
	}

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
		return object.NewError(fault.Import, node.Token, "module `%s` is still being imported, so it cannot be bound to a name yet", node.Path.Value).
			WithHelp("this is a circular import: %s", importCycle(node.Path.Value))
	}

	name := moduleBindingName(node.Path.Value)

	if node.Alias != nil {
		name = node.Alias.Value
	}

	scope.Environment.Set(name, moduleValue(moduleScope))

	return nil
}

// moduleBindingName derives the name a bare `import "path/to/module"` binds
// its module to, the same way most languages take the last path segment:
// `import "math"` binds `math`, `import "utils/math"` also binds `math`.
func moduleBindingName(path string) string {
	return filepath.Base(path)
}

// moduleValue turns a finished module's top-level bindings into a value that
// Ghost code can hold and read properties off of. A Map already supports dot
// access and has its own methods (keys(), has(), ...), so a loaded module
// reuses it rather than introducing a namespace type solely for this.
func moduleValue(moduleScope *object.Scope) *object.Map {
	pairs := make(map[object.MapKey]object.MapPair)

	for name, value := range moduleScope.Environment.All() {
		key := &object.String{Value: name}
		pairs[key.MapKey()] = object.MapPair{Key: key, Value: value}
	}

	return &object.Map{Pairs: pairs}
}

// evaluateSchemeImport handles the whole-value form of a `scheme:name`
// import — `import name from "scheme:name"` and `import name as alias from
// "scheme:name"`. The module or function itself — the same value a bare
// `console` resolves to, for the standard library's own scheme — is bound
// directly, so dot access on a module (`math.pi`, `math.sqrt()`) works
// exactly as it does on the still-global modules; only how the name reaches
// the scope differs. Which scheme is named doesn't matter here: the standard
// library's `ghost:` and an embedder's own scheme (Lumen's `lumen:`, say)
// resolve through the exact same lookup.
func evaluateSchemeImport(scheme string, name string, node *ast.Import, importScope *object.Scope) object.Object {
	value, err := lookupSchemeBinding(node.Token, scheme, name)

	if err != nil {
		return err
	}

	binding := name

	if node.Alias != nil {
		binding = node.Alias.Value
	}

	importScope.Environment.Set(binding, value)

	return nil
}

// evaluateSchemeImportFrom handles `import { a, b as c } from "scheme:name"`
// and `import * from "scheme:name"`. A method pulls in the same
// *object.LibraryFunction a dotted call would use; a property is evaluated
// once, immediately, since there is no lazy-getter value to bind — `import {
// pi } from "ghost:math"` binds a plain number, the way `math.pi` would read
// as one.
func evaluateSchemeImportFrom(scheme string, name string, node *ast.ImportFrom, importScope *object.Scope) object.Object {
	module, err := lookupSchemeModule(node.Token, scheme, name)

	if err != nil {
		return err
	}

	if node.Everything {
		for methodName, function := range module.Methods {
			importScope.Environment.Set(methodName, function)
		}

		for propertyName, property := range module.Properties {
			value := unwrapCall(node.Token, property, nil, importScope)

			if isError(value) {
				return value
			}

			importScope.Environment.Set(propertyName, value)
		}

		return nil
	}

	for alias, identifier := range node.Identifiers {
		if function, ok := module.Methods[identifier.Value]; ok {
			importScope.Environment.Set(alias, function)

			continue
		}

		if property, ok := module.Properties[identifier.Value]; ok {
			value := unwrapCall(node.Token, property, nil, importScope)

			if isError(value) {
				return value
			}

			importScope.Environment.Set(alias, value)

			continue
		}

		raised := object.NewError(fault.Import, node.Token, "module `%s` does not define `%s`", name, identifier.Value)

		if suggestion, ok := nearestName(identifier.Value, schemeModuleExports(module)); ok {
			raised.WithHelp("did you mean `%s`?", suggestion)
		}

		return raised
	}

	return nil
}

// lookupSchemeModule resolves a `scheme:name` import path against that
// scheme's own registry, for the `import { a, b } from "scheme:name"` /
// `import * from "scheme:name"` forms, which need something with methods and
// properties to pull members out of. It works for any scheme an embedder has
// registered under (library.RegisterModuleForScheme) just as it does for the
// standard library's own `ghost:` — including one registered mid-script by
// `ghost.extend` (§9.12), since a plugin's module becomes importable the
// moment it registers, with no separate mechanism to keep in sync.
func lookupSchemeModule(tok token.Token, scheme string, name string) (*object.LibraryModule, *object.Error) {
	registry := library.Scheme(scheme)

	if module, ok := registry.Modules[name]; ok {
		return module, nil
	}

	if _, ok := registry.Functions[name]; ok {
		return nil, object.NewError(fault.Import, tok, "`%s` is a standalone function under `%s:`, not a module", name, scheme).
			WithHelp("import it directly: `import \"%s:%s\"`", scheme, name)
	}

	if !knownScheme(scheme) {
		return nil, unknownScheme(tok, scheme)
	}

	raised := object.NewError(fault.Import, tok, "no module named `%s` registered under `%s:`", name, scheme)

	if suggestion, ok := nearestName(name, schemeModuleNames(registry)); ok {
		raised.WithHelp("did you mean `import { ... } from \"%s:%s\"`?", scheme, suggestion)
	}

	return nil, raised
}

// lookupSchemeBinding resolves a `scheme:name` import path for the
// whole-value forms, `import name from "scheme:name"` and `import name as
// alias from "scheme:name"`, which bind the entry itself — a module (dot
// access works exactly as it does on `console`) or a standalone function
// (directly callable), whichever `name` turns out to be.
func lookupSchemeBinding(tok token.Token, scheme string, name string) (object.Object, *object.Error) {
	registry := library.Scheme(scheme)

	if module, ok := registry.Modules[name]; ok {
		return module, nil
	}

	if function, ok := registry.Functions[name]; ok {
		return function, nil
	}

	if !knownScheme(scheme) {
		return nil, unknownScheme(tok, scheme)
	}

	raised := object.NewError(fault.Import, tok, "no module or function named `%s` registered under `%s:`", name, scheme)

	if suggestion, ok := nearestName(name, append(schemeModuleNames(registry), schemeFunctionNames(registry)...)); ok {
		raised.WithHelp("did you mean `import \"%s:%s\"`?", scheme, suggestion)
	}

	return nil, raised
}

// knownScheme reports whether anything has ever registered a module or
// function under this scheme, telling "wrong name within a real scheme"
// apart from "the scheme prefix itself is wrong or not loaded yet" — the
// latter needs a different message, since the fix isn't a nearby name within
// the scheme but the scheme prefix itself (or, for an embedder's scheme not
// registered until some setup step runs, waiting until after that step).
func knownScheme(scheme string) bool {
	for _, known := range library.Schemes() {
		if known == scheme {
			return true
		}
	}

	return false
}

// unknownScheme reports an import naming a scheme prefix nothing has ever
// registered under, suggesting the nearest one that does exist.
func unknownScheme(tok token.Token, scheme string) *object.Error {
	raised := object.NewError(fault.Import, tok, "no scheme named `%s:` is registered", scheme)

	if suggestion, ok := nearestName(scheme, library.Schemes()); ok {
		raised.WithHelp("did you mean `%s:`?", suggestion)
	}

	return raised
}

// schemeModuleNames lists every module registered under a scheme, for
// suggesting the one a misspelled import path probably meant.
func schemeModuleNames(registry *library.Registry) []string {
	names := make([]string, 0, len(registry.Modules))

	for name := range registry.Modules {
		names = append(names, name)
	}

	return names
}

// schemeFunctionNames is schemeModuleNames for standalone functions.
func schemeFunctionNames(registry *library.Registry) []string {
	names := make([]string, 0, len(registry.Functions))

	for name := range registry.Functions {
		names = append(names, name)
	}

	return names
}

// schemeModuleExports lists the names a module actually offers — its methods
// and its properties together — for suggesting the one a misspelled `import
// { x } from "scheme:..."` probably meant.
func schemeModuleExports(module *object.LibraryModule) []string {
	names := make([]string, 0, len(module.Methods)+len(module.Properties))

	for name := range module.Methods {
		names = append(names, name)
	}

	for name := range module.Properties {
		names = append(names, name)
	}

	return names
}

func evaluateImportFrom(node *ast.ImportFrom, scope *object.Scope) object.Object {
	if scheme, name, ok := schemeImport(node.Path.Value); ok {
		return evaluateSchemeImportFrom(scheme, name, node, scope)
	}

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
