package environment

import "ghostlang.org/x/ghost/object"

type Environment struct {
	store map[string]object.Object
}

func NewEnvironment() *Environment {
	store := make(map[string]object.Object)

	return &Environment{store: store}
}

func (environment *Environment) Get(name string) (object.Object, bool) {
	object, ok := environment.store[name]

	return object, ok
}

func (environment *Environment) Set(name string, value object.Object) object.Object {
	environment.store[name] = value

	return value
}
