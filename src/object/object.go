package object

// Type is the type of the object given as a string.
type Type string

// Object is the interface for all object values.
type Object interface {
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
	visitBoolean(*Boolean)
	visitFunction(*Function)
	visitLibraryFunction(*LibraryFunction)
	visitLibraryModule(*LibraryModule)
	visitList(*List)
	visitMap(*Map)
	visitNull(*Null)
	visitNumber(*Number)
	visitString(*String)
}

type Mappable interface {
	MapKey() MapKey
}

type HasMethods interface {
	Method(method string, args []Object) (Object, bool)
}

type HasVisitor interface {
	Accept(Visitor)
}

type GoFunction func(args ...Object) Object

type ObjectMethod func(value interface{}, args ...Object) (Object, bool)
