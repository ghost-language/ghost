package library

import (
	"ghostlang.org/x/ghost/library/functions"
	"ghostlang.org/x/ghost/library/modules"
	"ghostlang.org/x/ghost/object"
)

var Functions = map[string]*object.LibraryFunction{}
var Modules = map[string]*object.LibraryModule{}

func init() {
	RegisterModule("console", modules.ConsoleMethods, modules.ConsoleProperties)
	RegisterModule("ghost", modules.GhostMethods, modules.GhostProperties)
	RegisterModule("http", modules.HttpMethods, modules.HttpProperties)
	RegisterModule("io", modules.IoMethods, modules.IoProperties)
	RegisterModule("math", modules.MathMethods, modules.MathProperties)
	RegisterModule("os", modules.OsMethods, modules.OsProperties)
	RegisterModule("random", modules.RandomMethods, modules.RandomProperties)
	RegisterModule("time", modules.TimeMethods, modules.TimeProperties)
	RegisterModule("window", modules.WindowMethods, modules.WindowProperties)
	RegisterModule("draw", modules.DrawMethods, modules.DrawProperties)
	RegisterModule("text", modules.TextMethods, modules.TextProperties)
	RegisterModule("keyboard", modules.KeyboardMethods, modules.KeyboardProperties)
	RegisterModule("color", modules.ColorMethods, modules.ColorProperties)

	RegisterFunction("print", functions.Print)
	RegisterFunction("type", functions.Type)
}

func RegisterFunction(name string, function object.GoFunction) {
	Functions[name] = &object.LibraryFunction{Name: name, Function: function}
}

func RegisterModule(name string, methods map[string]*object.LibraryFunction, properties map[string]*object.LibraryProperty) {
	Modules[name] = &object.LibraryModule{Name: name, Methods: methods, Properties: properties}
}
