package modules

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

var evaluate func(node ast.Node, env *object.Environment) object.Object

func RegisterMethod(module map[string]*object.LibraryFunction, method string, function object.GoFunction) {
	module[method] = &object.LibraryFunction{Name: method, Function: function}
}

func RegisterEvaluator(e func(node ast.Node, env *object.Environment) object.Object) {
	evaluate = e
}
