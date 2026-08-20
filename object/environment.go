package object

import (
	"io"
	"os"
)

// inlineCapacity is how many variables an environment stores inline, in the
// environment struct itself, before it falls back to a map.
//
// Every function call builds a new environment, and almost all of them hold
// only a handful of names: the parameters plus a few locals. A map is a poor
// fit for that. Allocating one, and growing its buckets as names are set,
// dominated the memory profile of call-heavy programs. Storing the first few
// names inline means an environment costs a single allocation and lookups are
// a short scan of adjacent memory rather than a hash. Environments that hold
// more names than this - module and class scopes, mainly - spill the remainder
// into the overflow map and keep their previous behavior.
const inlineCapacity = 4

type Environment struct {
	names    [inlineCapacity]string
	values   [inlineCapacity]Object
	count    int
	overflow map[string]Object

	outer     *Environment
	writer    io.Writer
	directory string
}

func NewEnvironment() *Environment {
	return &Environment{writer: os.Stdout}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	environment := NewEnvironment()
	environment.outer = outer
	environment.writer = outer.writer

	return environment
}

// local looks a name up in this environment only, ignoring the outer chain.
func (environment *Environment) local(name string) (Object, bool) {
	for index := 0; index < environment.count; index++ {
		if environment.names[index] == name {
			return environment.values[index], true
		}
	}

	if environment.overflow != nil {
		value, ok := environment.overflow[name]

		return value, ok
	}

	return nil, false
}

// All returns a copy of the names bound in this environment.
func (environment *Environment) All() map[string]Object {
	all := make(map[string]Object, environment.count+len(environment.overflow))

	for index := 0; index < environment.count; index++ {
		all[environment.names[index]] = environment.values[index]
	}

	for name, value := range environment.overflow {
		all[name] = value
	}

	return all
}

func (environment *Environment) Has(name string) bool {
	_, ok := environment.local(name)

	if !ok && environment.outer != nil {
		_, ok = environment.outer.Get(name)
	}

	return ok
}

func (environment *Environment) Get(name string) (Object, bool) {
	object, ok := environment.local(name)

	if !ok && environment.outer != nil {
		object, ok = environment.outer.Get(name)
	}

	return object, ok
}

// HasLocal reports whether the name is bound in this environment itself,
// ignoring the enclosing chain. Instance fields are looked up this way so that
// a same-named variable in an enclosing scope cannot masquerade as a property.
func (environment *Environment) HasLocal(name string) bool {
	_, ok := environment.local(name)

	return ok
}

// GetLocal reads a name bound in this environment itself, ignoring the
// enclosing chain.
func (environment *Environment) GetLocal(name string) (Object, bool) {
	return environment.local(name)
}

// Set binds a name in this environment, replacing any existing binding here.
// It never walks the outer chain.
func (environment *Environment) Set(name string, value Object) Object {
	for index := 0; index < environment.count; index++ {
		if environment.names[index] == name {
			environment.values[index] = value

			return value
		}
	}

	if environment.overflow != nil {
		if _, ok := environment.overflow[name]; ok {
			environment.overflow[name] = value

			return value
		}
	}

	if environment.count < inlineCapacity {
		environment.names[environment.count] = name
		environment.values[environment.count] = value
		environment.count++

		return value
	}

	if environment.overflow == nil {
		environment.overflow = make(map[string]Object)
	}

	environment.overflow[name] = value

	return value
}

func (environment *Environment) Delete(name string) {
	for index := 0; index < environment.count; index++ {
		if environment.names[index] != name {
			continue
		}

		// Order is not meaningful here, so close the gap with the last entry
		// and clear it so the removed value can be collected.
		last := environment.count - 1

		environment.names[index] = environment.names[last]
		environment.values[index] = environment.values[last]
		environment.names[last] = ""
		environment.values[last] = nil
		environment.count--

		return
	}

	delete(environment.overflow, name)
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

// create a new function "Call" that can be used to call a function within the environment.
func (environment *Environment) Call(function string, args []Object, writer io.Writer) Object {
	if object, ok := environment.Get(function); ok {
		if function, ok := object.(*Function); ok {
			return function.Evaluate(args, writer)
		}
	}

	return NewError("function not found: %s", function)
}
