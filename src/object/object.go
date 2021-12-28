package object

// Type is the type of the object given as a string.
type Type string

// Object is the interface for all object values.
type Object interface {
	Visitable
	Type() Type
	String() string
}

type MapKey struct {
	Type  Type
	Value uint64
}

type Visitor interface {
	visitBoolean(*Boolean)
	visitFunction(*Function)
	visitLibraryFunction(*LibraryFunction)
	visitList(*List)
	visitNull(*Null)
	visitNumber(*Number)
	visitString(*String)
}

type Visitable interface {
	Accept(Visitor)
}

type GoFunction func(args ...Object) Object
