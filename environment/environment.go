package environment

import (
	"fmt"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

// Environment stores the bindings that associate variables to values.
type Environment struct {
	values    map[string]object.Object
	enclosing *Environment
}

// New creates a new instance of Environment.
func New() *Environment {
	return &Environment{values: make(map[string]object.Object), enclosing: nil}
}

// Extend an existing environment.
func Extend(env *Environment) *Environment {
	return &Environment{values: make(map[string]object.Object), enclosing: env}
}

// Define binds a new value to the environment with the given name.
func (e *Environment) Define(name string, value object.Object) {
	e.values[name] = value
}

// Get fetches the variable with the given name from the environment.
func (e *Environment) Get(name token.Token) (object.Object, error) {
	result, exists := e.values[name.Lexeme]

	if exists {
		return result, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	return nil, fmt.Errorf("Undefined variable '%v'", name.Lexeme)
}
