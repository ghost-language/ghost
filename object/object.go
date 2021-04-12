package object

// Type is the type of the token given as a string
type Type string

// Object is the interface for all object values.
type Object interface {
	Type() Type
	String() string
}

type GhostFunction func(env *Environment, args ...Object) Object