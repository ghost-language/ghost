package library

import (
	"sort"

	"ghostlang.org/x/ghost/library/functions"
	"ghostlang.org/x/ghost/library/modules"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/optimizer"
)

// StandardScheme is the scheme name the standard library itself registers
// under — `ghost:math`, `ghost:file`, and so on (§8.9 in SPEC.md).
const StandardScheme = "ghost"

// Functions and Modules are StandardScheme's own registries. They stay
// package-level vars, rather than folding into Registry/schemes below, so
// existing code that reads them directly (evaluator/errors.go's typo
// suggestions, the tests) keeps working unchanged; RegisterFunction and
// RegisterModule write into these same maps, and schemes[StandardScheme]
// below is set up to alias them rather than duplicate them.
var Functions = map[string]*object.LibraryFunction{}
var Modules = map[string]*object.LibraryModule{}

// globalModules and globalFunctions name the only library bindings reachable
// without an import — everything else in the standard library (math, file,
// path, os, random, json, http, date, ghost), and everything any embedder
// registers under any scheme, has to be pulled in with `import ... from
// "scheme:name"` (see evaluator/import.go). console and type are common
// enough, and small enough in surface, to earn the same standing as a
// keyword.
var globalModules = map[string]bool{"console": true}
var globalFunctions = map[string]bool{"type": true}

// Registry is one scheme's own modules and functions — the standard
// library's ("ghost"), or one an embedding host claims for itself.
type Registry struct {
	Modules   map[string]*object.LibraryModule
	Functions map[string]*object.LibraryFunction
}

// schemes holds every import scheme registered so far, keyed by scheme name.
// An embedding host — Lumen, say — claims one with RegisterModuleForScheme/
// RegisterFunctionForScheme, and a script then reaches it the same way it
// reaches the standard library, just with its own prefix: `import font from
// "lumen:font"` instead of `import math from "ghost:math"`. There is nothing
// standard-library-specific about the mechanism; StandardScheme is simply
// the one scheme Ghost itself pre-registers into.
var schemes = map[string]*Registry{
	StandardScheme: {Modules: Modules, Functions: Functions},
}

func init() {
	RegisterModule("console", modules.ConsoleMethods, modules.ConsoleProperties)
	RegisterModule("date", modules.DateMethods, modules.DateProperties)
	RegisterModule("ghost", modules.GhostMethods, modules.GhostProperties)
	RegisterModule("http", modules.HttpMethods, modules.HttpProperties)
	RegisterModule("file", modules.FileMethods, modules.FileProperties)
	RegisterModule("path", modules.PathMethods, modules.PathProperties)
	RegisterModule("math", modules.MathMethods, modules.MathProperties)
	RegisterModule("os", modules.OsMethods, modules.OsProperties)
	RegisterModule("random", modules.RandomMethods, modules.RandomProperties)
	RegisterModule("json", modules.JsonMethods, modules.JsonProperties)

	RegisterFunction("type", functions.Type)

	optimizer.SetGlobalResolver(IsGlobal)
}

// IsGlobal reports whether a name is reachable without an import: console and
// type, and nothing else. The optimizer uses it to classify identifiers ahead
// of evaluation.
func IsGlobal(name string) bool {
	return globalModules[name] || globalFunctions[name]
}

// GlobalModule reports the module a bare name resolves to without an import,
// if any. Evaluating a bare identifier consults this rather than Modules
// directly, so a module registered for import (by the standard library or by
// an embedder, under any scheme) does not leak into ordinary scope just by
// existing.
func GlobalModule(name string) (*object.LibraryModule, bool) {
	if !globalModules[name] {
		return nil, false
	}

	module, ok := Modules[name]

	return module, ok
}

// GlobalFunction is GlobalModule for the (currently sole) global function.
func GlobalFunction(name string) (*object.LibraryFunction, bool) {
	if !globalFunctions[name] {
		return nil, false
	}

	function, ok := Functions[name]

	return function, ok
}

// RegisterFunction and RegisterModule register into the standard library's
// own scheme (`ghost:name`) — the same signatures the embedding API has
// always had, so existing embedding code keeps compiling and keeps reaching
// what it registers at the same import path it always has.
func RegisterFunction(name string, function object.GoFunction) {
	RegisterFunctionForScheme(StandardScheme, name, function)
}

func RegisterModule(name string, methods map[string]*object.LibraryFunction, properties map[string]*object.LibraryProperty) {
	RegisterModuleForScheme(StandardScheme, name, methods, properties)
}

// RegisterFunctionForScheme and RegisterModuleForScheme register a function
// or module under a scheme of the caller's own choosing, rather than the
// standard library's. This is how a Go program embedding Ghost — Lumen,
// say — hands its own scripts a namespace that reads as its own:
//
//	ghost.RegisterModuleForScheme("lumen", "font", methods, properties)
//	// -> import font from "lumen:font"
//	// -> import { load } from "lumen:font"
//
// There is one registry and one `scheme:name` import mechanism underneath
// every scheme, StandardScheme included, so an embedder gets exactly the
// same resolution, typo suggestions, and "not registered yet" behavior the
// standard library gets — nothing about `ghost:` is special-cased.
func RegisterFunctionForScheme(scheme string, name string, function object.GoFunction) {
	registry(scheme).Functions[name] = &object.LibraryFunction{Name: name, Function: function}
}

func RegisterModuleForScheme(scheme string, name string, methods map[string]*object.LibraryFunction, properties map[string]*object.LibraryProperty) {
	registry(scheme).Modules[name] = &object.LibraryModule{Name: name, Methods: methods, Properties: properties}
}

// Scheme returns the registry a `scheme:name` import resolves against,
// creating an empty one on first use — importing from a scheme nothing has
// registered under yet fails the same way importing an unknown name within a
// real scheme does (a lookup that finds nothing), not with a different error
// shape. Callers that need to tell "nobody has ever registered under this
// scheme" apart from "this scheme exists but doesn't have that name" should
// check Schemes() first.
func Scheme(scheme string) *Registry {
	return registry(scheme)
}

// Schemes lists every scheme with at least one module or function registered
// — StandardScheme always among them — for suggesting the right one when a
// script's `import` names a scheme prefix nothing has claimed.
func Schemes() []string {
	names := make([]string, 0, len(schemes))

	for name, reg := range schemes {
		if len(reg.Modules) > 0 || len(reg.Functions) > 0 {
			names = append(names, name)
		}
	}

	sort.Strings(names)

	return names
}

// LocateModule reports every registered scheme that has a module by this
// name — usually zero or one, but nothing stops two schemes from both
// registering, say, "font". Used to suggest the right `import` for a bare
// name that used to work as a global, or that only makes sense as a module.
func LocateModule(name string) []string {
	return locate(name, func(registry *Registry) bool {
		_, ok := registry.Modules[name]

		return ok
	})
}

// LocateFunction is LocateModule for standalone functions.
func LocateFunction(name string) []string {
	return locate(name, func(registry *Registry) bool {
		_, ok := registry.Functions[name]

		return ok
	})
}

// Locate is LocateModule and LocateFunction combined — every scheme
// registering either a module or a function by this name, since a caller
// asking "where would I import this from" rarely cares which kind it is:
// the whole-value form of an import (`import name from "scheme:name"`)
// binds either one identically.
func Locate(name string) []string {
	return locate(name, func(registry *Registry) bool {
		if _, ok := registry.Modules[name]; ok {
			return true
		}

		_, ok := registry.Functions[name]

		return ok
	})
}

func locate(name string, has func(*Registry) bool) []string {
	found := make([]string, 0, 1)

	for scheme, registry := range schemes {
		if has(registry) {
			found = append(found, scheme)
		}
	}

	sort.Strings(found)

	return found
}

func registry(scheme string) *Registry {
	reg, ok := schemes[scheme]

	if !ok {
		reg = &Registry{Modules: map[string]*object.LibraryModule{}, Functions: map[string]*object.LibraryFunction{}}
		schemes[scheme] = reg
	}

	return reg
}
