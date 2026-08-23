package library

import (
	"ghostlang.org/x/ghost/library/functions"
	"ghostlang.org/x/ghost/library/modules"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/optimizer"
)

var Functions = map[string]*object.LibraryFunction{}
var Modules = map[string]*object.LibraryModule{}

func init() {
	RegisterModule("console", modules.ConsoleMethods, modules.ConsoleProperties)
	RegisterModule("date", modules.DateMethods, modules.DateProperties)
	RegisterModule("ghost", modules.GhostMethods, modules.GhostProperties)
	RegisterModule("http", modules.HttpMethods, modules.HttpProperties)
	RegisterModule("io", modules.IoMethods, modules.IoProperties)
	RegisterModule("math", modules.MathMethods, modules.MathProperties)
	RegisterModule("os", modules.OsMethods, modules.OsProperties)
	RegisterModule("random", modules.RandomMethods, modules.RandomProperties)
	RegisterModule("json", modules.JsonMethods, modules.JsonProperties)

	RegisterFunction("print", functions.Print)
	RegisterFunction("type", functions.Type)

	optimizer.SetGlobalResolver(IsGlobal)
}

// IsGlobal reports whether a name refers to a registered library module or
// function. The optimizer uses it to classify identifiers ahead of evaluation.
func IsGlobal(name string) bool {
	if _, ok := Modules[name]; ok {
		return true
	}

	_, ok := Functions[name]

	return ok
}

func RegisterFunction(name string, function object.GoFunction) {
	Functions[name] = &object.LibraryFunction{Name: name, Function: function}
}

func RegisterModule(name string, methods map[string]*object.LibraryFunction, properties map[string]*object.LibraryProperty) {
	Modules[name] = &object.LibraryModule{Name: name, Methods: methods, Properties: properties}
}
