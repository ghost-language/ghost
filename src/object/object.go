package object

import (
	"ghostlang.org/x/ghost/ast"
)

// Type is the type of the object given as a string.
type Type string

// Object is the interface for all object values.
type Object interface {
	HasMethods
	Type() Type
	String() string
}

type MapKey struct {
	Type  Type
	Value uint64
}

type Mappable interface {
	MapKey() MapKey
}

type HasMethods interface {
	Method(method string, args []Object) (Object, bool)
}

type GoFunction func(args ...Object) Object

type ObjectMethod func(value interface{}, args ...Object) (Object, bool)

var evaluator func(node ast.Node, env *Environment) (Object, bool)

func SetEvaluator(e func(node ast.Node, env *Environment) (Object, bool)) {
	evaluator = e
}
