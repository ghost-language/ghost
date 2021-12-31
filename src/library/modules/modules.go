package modules

import "ghostlang.org/x/ghost/object"

func RegisterMethod(module []*object.LibraryFunction, method string, function object.GoFunction) []*object.LibraryFunction {
	module = append(module, &object.LibraryFunction{Name: method, Function: function})

	return module
}
