package modules

import "ghostlang.org/x/ghost/object"

func RegisterMethod(module map[string]*object.LibraryFunction, method string, function object.GoFunction) {
	module[method] = &object.LibraryFunction{Name: method, Function: function}
}
