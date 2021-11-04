package library

import (
	"ghostlang.org/x/ghost/library/functions"
	"ghostlang.org/x/ghost/object"
)

var Functions = map[string]*object.LibraryFunction{}

func init() {
	RegisterFunction("print", functions.Print)
	RegisterFunction("sleep", functions.Sleep)
	RegisterFunction("type", functions.Type)
}

func RegisterFunction(name string, function object.GoFunction) {
	Functions[name] = &object.LibraryFunction{Name: name, Function: function}
}
