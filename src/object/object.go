package object

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

var evaluator func(node ast.Node, env *Environment) Object

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

type GoFunction func(env *Environment, tok token.Token, args ...Object) Object
type ObjectMethod func(value interface{}, args ...Object) (Object, bool)

func RegisterEvaluator(e func(node ast.Node, env *Environment) Object) {
	evaluator = e
}
