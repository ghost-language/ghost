package library

import (
	"ghostlang.org/x/ghost/library/functions"
	"ghostlang.org/x/ghost/library/modules"
	"ghostlang.org/x/ghost/object"
)

var Functions = map[string]*object.LibraryFunction{}
var Modules = map[string]*object.LibraryModule{}

func init() {
	RegisterModule("math", modules.Math)

	RegisterFunction("print", functions.Print)
	RegisterFunction("readLine", functions.ReadLine)
	RegisterFunction("sleep", functions.Sleep)
	RegisterFunction("type", functions.Type)
}

func RegisterFunction(name string, function object.GoFunction) {
	Functions[name] = &object.LibraryFunction{Name: name, Function: function}
}

func RegisterModule(name string, methods map[string]*object.LibraryFunction) {
	Modules[name] = &object.LibraryModule{Name: name, Methods: methods}
}
