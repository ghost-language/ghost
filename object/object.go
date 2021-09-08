package object

import "ghostlang.org/x/ghost/token"

// Type is the type of the token given as a string
type Type string

// Object is the interface for all object values.
type Object interface {
	Type() Type
	String() string
}

type PropertyAccessor interface {
	Get(name token.Token) Object
}

type GhostFunction func(env *Environment, args ...Object) Object