package object

import (
	"fmt"

	"ghostlang.org/x/ghost/token"
)

const INSTANCE = "INSTANCE"

// Instance objects consist of a body and an environment.
type Instance struct {
	Class       *Class
	Environment *Environment
}

// String represents the instance object's value as a string.
func (instance *Instance) String() string {
	return fmt.Sprintf("class instance %s", instance.Class.Name.Value)
}

// Type returns the instance object type.
func (instance *Instance) Type() Type {
	return INSTANCE
}

// Method defines the set of methods available on instance objects.
func (instance *Instance) Method(method string, args []Object) (Object, bool) {
	return nil, false
}

func (instance *Instance) Call(name string, arguments []Object, tok token.Token) Object {
	if function, ok := instance.Environment.Get(name); ok {
		if method, ok := function.(*Function); ok {
			functionEnvironment := createMethodEnvironment(method, arguments)
			functionScope := &Scope{Self: instance, Environment: functionEnvironment}

			return evaluator(method.Body, functionScope)
		}
	}

	return NewError("%d:%d: runtime error: unknown method '%s' on class %s", tok.Line, tok.Column, name, instance.Class.Name.Value)
}

func createMethodEnvironment(method *Function, arguments []Object) *Environment {
	env := NewEnclosedEnvironment(method.Scope.Environment)

	for key, val := range method.Defaults {
		env.Set(key, evaluator(val, method.Scope))
	}

	for index, parameter := range method.Parameters {
		if index < len(arguments) {
			env.Set(parameter.Value, arguments[index])
		}
	}

	return env
}
