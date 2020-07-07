package object

// NewEnvironment constructs a new Environment object to hold bindings
// of identifiers to their names
func NewEnvironment() *Environment {
	s := make(map[string]Object)

	return &Environment{store: s}
}

// Environment is an object that holds a mapping of names to bound objects
type Environment struct {
	store map[string]Object
}

// Get returns the object bound by name
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]

	return obj, ok
}

// Set stores the object with the given name
func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value

	return value
}
