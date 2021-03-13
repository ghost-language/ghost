package environment

import (
	"bytes"
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

func (e *Environment) All() map[string]object.Object {
	return e.values
}

// Set binds a new value to the environment with the given name.
func (e *Environment) Set(name string, value object.Object) object.Object {
	e.values[name] = value

	return value
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

func (e *Environment) String() string {
	var values bytes.Buffer

	for key, value := range e.values {
		values.WriteString(fmt.Sprintf("%s: %s\n", key, value.String()))
	}

	return values.String()
}

// SetWriter ...
func (e *Environment) SetWriter(writer io.Writer) {
	e.writer = writer
}

// GetWriter ...
func (e *Environment) GetWriter() io.Writer {
	return e.writer
}
