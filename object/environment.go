package object

import (
	"io"
	"os"

	"ghostlang.org/x/ghost/fault"
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

	// captured records that a value outliving this environment holds a
	// reference to it — a closure, a class, or a trait created while it was in
	// scope. Such an environment can never be reused, because whatever
	// captured it will read from it later. Capture marks the whole enclosing
	// chain, since a closure reads through `outer` as well.
	captured bool

	// freeChild holds at most one finished block scope, so a loop body or an
	// `if` inside a hot loop reuses a single scope and environment rather than
	// allocating a pair per execution. A block's scope is only ever reached
	// from the one goroutine executing that block: every concurrent entry into
	// Ghost code (an `http.handle` callback, an embedder's Call) runs in a
	// function frame of its own, and a block's environment is a child of that
	// frame rather than of anything shared.
	freeChild *Scope
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

// Capture marks this environment, and every environment enclosing it, as held
// by a value that outlives the block it belongs to. Creating a closure, a
// class, or a trait captures the scope it was created in, and reading a name
// through that value later walks the whole chain — so none of it may be
// reused.
func (environment *Environment) Capture() {
	for current := environment; current != nil; current = current.outer {
		if current.captured {
			return
		}

		current.captured = true
	}
}

// clear empties an environment so it can be handed to another block, keeping
// only what makes it that block's environment rather than any other: where it
// is chained, and where it writes.
func (environment *Environment) clear() {
	for index := 0; index < environment.count; index++ {
		environment.names[index] = ""
		environment.values[index] = nil
	}

	environment.count = 0
	environment.extra = nil
	environment.overflow = nil
	environment.captured = false
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

// rebind replaces the value under a name already bound in this environment,
// reporting whether there was one. It never creates a binding, which is what
// separates it from Set and what lets Assign walk outward without declaring a
// name in every scope it passes through.
func (environment *Environment) rebind(name string, value Object) bool {
	for index := 0; index < environment.count; index++ {
		if environment.names[index] == name {
			environment.values[index] = value

			return true
		}
	}

	for index := range environment.extra {
		if environment.extra[index].name == name {
			environment.extra[index].value = value

			return true
		}
	}

	if environment.overflow != nil {
		if _, ok := environment.overflow[name]; ok {
			environment.overflow[name] = value

			return true
		}
	}

	return false
}

// Assign rebinds a name wherever it is already bound, walking outward through
// the enclosing chain the same way Get does, and reports whether it found one.
// A name bound nowhere is left alone for the caller to declare, so assignment
// reaches an existing outer variable (§13.13) without silently creating outer
// bindings for names that are genuinely new here.
func (environment *Environment) Assign(name string, value Object) bool {
	for current := environment; current != nil; current = current.outer {
		if current.rebind(name, value) {
			return true
		}
	}

	return false
}

// Set binds a name in this environment, replacing any existing binding here.
// It never walks the outer chain — Assign is the one that does.
func (environment *Environment) Set(name string, value Object) Object {
	if environment.rebind(name, value) {
		return value
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

// Call invokes a named function in this environment. It is how a Go program
// embedding Ghost reaches into a script, so there is no token to point at and
// the fault it raises carries no position.
func (environment *Environment) Call(function string, args []Object, writer io.Writer) Object {
	value, ok := environment.Get(function)

	if !ok {
		return NewErrorFrom(fault.New(fault.Name, "`%s` is not defined", function))
	}

	callable, ok := value.(*Function)

	if !ok {
		return NewErrorFrom(fault.New(fault.Type, "`%s` is a %s, which cannot be called", function, TypeName(value)))
	}

	return callable.Evaluate(args, writer)
}
