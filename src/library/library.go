package library

import (
	"ghostlang.org/x/ghost/library/functions"
	"ghostlang.org/x/ghost/library/modules"
	"ghostlang.org/x/ghost/object"
)

var Functions = map[string]*object.LibraryFunction{}
var Modules = map[string]*object.LibraryModule{}

func init() {
	RegisterModule("console", modules.Console)
	RegisterModule("ghost", modules.Ghost)
	RegisterModule("http", modules.Http)
	RegisterModule("math", modules.Math)
	RegisterModule("os", modules.Os)

	RegisterFunction("print", functions.Print)
	RegisterFunction("type", functions.Type)
}

func RegisterFunction(name string, function object.GoFunction) {
	Functions[name] = &object.LibraryFunction{Name: name, Function: function}
}

func RegisterModule(name string, methods map[string]*object.LibraryFunction) {
	Modules[name] = &object.LibraryModule{Name: name, Methods: methods}
}
