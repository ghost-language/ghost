package object

import (
	"io"
	"os"
)

type Environment struct {
	store     map[string]Object
	outer     *Environment
	writer    io.Writer
	directory string
}

func NewEnvironment() *Environment {
	store := make(map[string]Object)

	return &Environment{store: store, writer: os.Stdout}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	environment := NewEnvironment()
	environment.outer = outer
	environment.writer = outer.writer

	return environment
}

func (environment *Environment) All() map[string]Object {
	return environment.store
}

func (environment *Environment) Has(name string) bool {
	_, ok := environment.store[name]

	if !ok && environment.outer != nil {
		_, ok = environment.outer.Get(name)
	}

	return ok
}

func (environment *Environment) Get(name string) (Object, bool) {
	object, ok := environment.store[name]

	if !ok && environment.outer != nil {
		object, ok = environment.outer.Get(name)
	}

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

func (environment *Environment) SetDirectory(directory string) {
	environment.directory = directory
}

func (environment *Environment) GetDirectory() string {
	directory := environment.directory

	if directory == "" && environment.outer != nil {
		directory = environment.outer.GetDirectory()
	}

	return directory
}
