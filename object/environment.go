package object

import (
	"io"
	"os"
)

// Environments store their bindings in three tiers, in lookup order: a fixed
// array inside the struct, a grown slice, and finally a map.
//
// The tiers exist because the two shapes of scope in a Ghost program want
// opposite things. A call frame holds a couple of parameters and is created
// again on every call, so it wants to be small and to cost a single allocation;
// the top level of a script holds every global and is read constantly from
// inside hot loops, so it wants lookups to stay cheap well past a couple of
// names. Sizing one fixed array for both made calls slower when it was large
// and script globals slower when it was small.
//
// Scanning beats hashing at these sizes by more than it looks like it should.
// A name being looked up and the name stored in the environment usually come
// from the same identifier in the AST, so the string comparison settles on
// equal pointers without examining any characters.
//
// Past scanLimit a scan really does stop paying for itself, so everything moves
// into the map and the scan tiers are abandoned. Splitting names across a scan
// tier and a map would be worse than either: every miss would pay a full scan
// before the map lookup.
const (
	inlineCapacity = 4
	scanLimit      = 16
)

// binding is one name/value pair in an environment's second storage tier.
type binding struct {
	name  string
	value Object
}

type Environment struct {
	// Tier one: inline, so an environment this small costs no allocation
	// beyond the struct itself.
	names  [inlineCapacity]string
	values [inlineCapacity]Object
	count  int

	// Tier two: allocated only for environments that outgrow the inline array.
	// One slice of pairs rather than parallel slices, because the header sits
	// in every environment, including the call frames that never grow one.
	extra []binding

	// Tier three: for environments larger than scanLimit. Once this exists the
	// scan tiers are empty.
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

	for index := range environment.extra {
		if environment.extra[index].name == name {
			return environment.extra[index].value, true
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
	all := make(map[string]Object, environment.count+len(environment.extra)+len(environment.overflow))

	for index := 0; index < environment.count; index++ {
		all[environment.names[index]] = environment.values[index]
	}

	for index := range environment.extra {
		all[environment.extra[index].name] = environment.extra[index].value
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

	for index := range environment.extra {
		if environment.extra[index].name == name {
			environment.extra[index].value = value

			return value
		}
	}

	if environment.overflow != nil {
		environment.overflow[name] = value

		return value
	}

	if environment.count < inlineCapacity {
		environment.names[environment.count] = name
		environment.values[environment.count] = value
		environment.count++

		return value
	}

	if environment.count+len(environment.extra) < scanLimit {
		environment.extra = append(environment.extra, binding{name: name, value: value})

		return value
	}

	// Too large to scan. Move every binding into the map and abandon the scan
	// tiers so that lookups cost one map access rather than a scan and a map
	// access.
	environment.overflow = make(map[string]Object, scanLimit*2)

	for index := 0; index < environment.count; index++ {
		environment.overflow[environment.names[index]] = environment.values[index]
		environment.names[index] = ""
		environment.values[index] = nil
	}

	for index := range environment.extra {
		environment.overflow[environment.extra[index].name] = environment.extra[index].value
	}

	environment.count = 0
	environment.extra = nil

	environment.overflow[name] = value

	return value
}

func (environment *Environment) Delete(name string) {
	// Order is not meaningful in either scan tier, so removals close the gap
	// with the last entry and clear it so the removed value can be collected.
	for index := 0; index < environment.count; index++ {
		if environment.names[index] != name {
			continue
		}

		last := environment.count - 1

		environment.names[index] = environment.names[last]
		environment.values[index] = environment.values[last]
		environment.names[last] = ""
		environment.values[last] = nil
		environment.count--

		return
	}

	for index := range environment.extra {
		if environment.extra[index].name != name {
			continue
		}

		last := len(environment.extra) - 1

		environment.extra[index] = environment.extra[last]
		environment.extra[last] = binding{}
		environment.extra = environment.extra[:last]

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
