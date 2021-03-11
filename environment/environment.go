package environment

import (
	"fmt"
	"io"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

// Environment stores the bindings that associate variables to values.
type Environment struct {
	values    map[string]object.Object
	enclosing *Environment
	writer    io.Writer
}

// New creates a new instance of Environment.
func New() *Environment {
	return &Environment{values: make(map[string]object.Object), enclosing: nil}
}

// Extend an existing environment.
func Extend(env *Environment) *Environment {
	return &Environment{values: make(map[string]object.Object), enclosing: env, writer: env.writer}
}

// Set binds a new value to the environment with the given name.
func (e *Environment) Set(name string, value object.Object) {
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

// SetWriter ...
func (e *Environment) SetWriter(writer io.Writer) {
	e.writer = writer
}

// GetWriter ...
func (e *Environment) GetWriter() io.Writer {
	return e.writer
}
