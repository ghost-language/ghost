package object

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/contract"
)

// Type is the type of the object given as a string.
type Type string

// Object is the interface for all object values.
type Object interface {
	contract.Object
	HasMethods
	HasVisitor
	Type() Type
	String() string
}

type MapKey struct {
	Type  Type
	Value uint64
}

type Visitor interface {
	VisitBoolean(*Boolean) (Object, bool)
	VisitFunction(*Function) (Object, bool)
	VisitLibraryFunction(*LibraryFunction) (Object, bool)
	VisitLibraryModule(*LibraryModule) (Object, bool)
	VisitList(*List) (Object, bool)
	VisitMap(*Map) (Object, bool)
	VisitNull(*Null) (Object, bool)
	VisitNumber(*Number) (Object, bool)
	VisitString(*String) (Object, bool)
}

type Mappable interface {
	MapKey() MapKey
}

type HasMethods interface {
	Method(method string, args []Object) (Object, bool)
}

type HasVisitor interface {
	Accept(Visitor) (Object, bool)
}

type GoFunction func(args ...Object) Object

type ObjectMethod func(value interface{}, args ...Object) (Object, bool)

var evaluator func(node ast.Node, env *Environment) (Object, bool)

func SetEvaluator(e func(node ast.Node, env *Environment) (Object, bool)) {
	evaluator = e
}
