package modules

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

var evaluate func(node ast.Node, scope *object.Scope) object.Object

func RegisterMethod(methods map[string]*object.LibraryFunction, name string, function object.GoFunction) {
	methods[name] = &object.LibraryFunction{Name: name, Function: function}
}

func RegisterProperty(properties map[string]*object.LibraryProperty, name string, property object.GoProperty) {
	properties[name] = &object.LibraryProperty{Name: name, Property: property}
}

func RegisterEvaluator(e func(node ast.Node, scope *object.Scope) object.Object) {
	evaluate = e
}
