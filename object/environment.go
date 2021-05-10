// Package object describes a named language entity such as
// a type, variable, function, or literal. All objects
// implement the Object interface.
package object

import (
	"io"
	"os"
	"unicode"
)

// Environment is an object that holds a mapping of names to bound objects
type Environment struct {
	store map[string]Object
	outer *Environment
	writer io.Writer
	directory string
}

// NewEnvironment constructs a new Environment object to hold bindings
// of identifiers to their names
func NewEnvironment() *Environment {
	s := make(map[string]Object)

	return &Environment{store: s, writer: os.Stdout, outer: nil}
}

// NewEnclosedEnvironment constructs a new Environment, extending off
// a pre-existing environment.
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	env.writer = outer.writer

	return env
}

// All returns all stored identifiers.
func (e *Environment) All() map[string]Object {
	return e.store
}

// Get returns the object bound by name
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]

	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

// Set stores the object with the given name
func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value

	return value
}

// Delete removes the object with the given name
func (e *Environment) Delete(name string) {
	delete(e.store, name)
}

func (e *Environment) Exported() *Map {
	pairs := make(map[MapKey]MapPair)

	for k, v := range e.store {
		// Replace this with checking for "Import" token
		if unicode.IsUpper(rune(k[0])) {
			s := &String{Value: k}
			pairs[s.MapKey()] = MapPair{Key: s, Value: v}
		}
	}

	return &Map{Pairs: pairs}
}

func (e *Environment) SetWriter(writer io.Writer) {
	e.writer = writer
}

func (e *Environment) GetWriter() io.Writer {
	return e.writer
}

func (e *Environment) SetDirectory(directory string) {
	e.directory = directory
}

func (e *Environment) GetDirectory() string {
	return e.directory
}