package builtins

import (
	"ghostlang.org/x/ghost/object"
)

// BuiltinFunctions stores a map of all registered native functions.
var BuiltinFunctions = map[string]*object.Builtin{}

// BuiltinModules stores a map of all registered native functions.
var BuiltinModules = map[string][]*object.Builtin{}

// RegisterFunction registers a new native function with Ghost.
func RegisterFunction(name string, function object.BuiltinFunction) {
	BuiltinFunctions[name] = &object.Builtin{Fn: function}
}

// RegisterModuleFunction registers a new native module with Ghost.
func RegisterModuleFunction(module string, name string, function object.BuiltinFunction) {
	BuiltinModules[module] = append(BuiltinModules[module], &object.Builtin{Fn: function})
}
