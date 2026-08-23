package library

import (
	"ghostlang.org/x/ghost/library/functions"
	"ghostlang.org/x/ghost/library/modules"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/optimizer"
)

var Functions = map[string]*object.LibraryFunction{}
var Modules = map[string]*object.LibraryModule{}

// globalModules and globalFunctions name the only library bindings reachable
// without an import — everything else in the standard library (math, file,
// path, os, random, json, http, date, ghost, and anything an embedder
// registers) has to be pulled in with `import ... from "ghost:name"` (see
// evaluator/import.go). console and type are common enough, and small enough
// in surface, to earn the same standing as a keyword.
var globalModules = map[string]bool{"console": true}
var globalFunctions = map[string]bool{"type": true}

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
// an embedder) does not leak into ordinary scope just by existing.
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

func RegisterFunction(name string, function object.GoFunction) {
	Functions[name] = &object.LibraryFunction{Name: name, Function: function}
}

func RegisterModule(name string, methods map[string]*object.LibraryFunction, properties map[string]*object.LibraryProperty) {
	Modules[name] = &object.LibraryModule{Name: name, Methods: methods, Properties: properties}
}
