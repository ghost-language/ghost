package builtins

import (
	"ghostlang.org/x/ghost/object"
)

// Builtins stores a map of all registered native functions.
var Builtins = map[string]*object.Builtin{}

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// RegisterFunction registers a new native function with Ghost.
func RegisterFunction(name string, function object.BuiltinFunction) {
	Builtins[name] = &object.Builtin{Fn: function}
}
