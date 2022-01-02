package object

import (
	"io"
	"os"
)

type Environment struct {
	store  map[string]Object
	writer io.Writer
}

func NewEnvironment() *Environment {
	store := make(map[string]Object)

	return &Environment{store: store, writer: os.Stdout}
}

func (environment *Environment) Get(name string) (Object, bool) {
	object, ok := environment.store[name]

	return object, ok
}

func (environment *Environment) Set(name string, value Object) Object {
	environment.store[name] = value

	return value
}

func (environment *Environment) Delete(name string) {
	delete(environment.store, name)
}

func (environment *Environment) SetWriter(writer io.Writer) {
	environment.writer = writer
}

func (environment *Environment) GetWriter() io.Writer {
	return environment.writer
}
