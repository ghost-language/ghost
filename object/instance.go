package object

import (
	"fmt"

	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/token"
)

// Instance objects consist of a body and an environment.
type Instance struct {
	Class *Class

	// Environment holds this instance's own fields. Lookups against it are
	// always local: it encloses the class environment only so that the output
	// writer and working directory are inherited.
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
func (instance *Instance) Method(method string, tok token.Token, args []Object) (Object, bool) {
	return nil, false
}

// LookupMember resolves a name against the instance's class chain, checking
// each class before the traits it uses and walking up to the superclass only
// after both. It returns the member along with the class that declared it; the
// caller records that class on the scope so `super` resolves relative to the
// declaration site rather than the receiver's class.
func (instance *Instance) LookupMember(name string) (Object, *Class, bool) {
	return LookupMember(instance.Class, name)
}

// LookupMember resolves a name starting at the given class and walking up the
// superclass chain. `super` uses it with the declaring class's superclass.
func LookupMember(start *Class, name string) (Object, *Class, bool) {
	for class := start; class != nil; class = class.Super {
		if member, ok := class.Environment.GetLocal(name); ok {
			return member, class, true
		}

		for _, trait := range class.Traits {
			if member, ok := trait.Environment.GetLocal(name); ok {
				return member, class, true
			}
		}
	}

	return nil, nil, false
}

// Call invokes a method on the instance by name. It is the entry point used by
// Go code embedding Ghost.
func (instance *Instance) Call(name string, arguments []Object, tok token.Token) Object {
	member, class, ok := instance.LookupMember(name)

	if ok {
		if method, ok := member.(*Function); ok {
			return instance.callMethod(method, class, arguments)
		}
	}

	if ok {
		return NewError(fault.Type, tok, "`%s.%s` is a %s, which cannot be called", instance.Class.Name.Value, name, TypeName(member))
	}

	return NewError(fault.Property, tok, "class `%s` has no method `%s`", instance.Class.Name.Value, name)
}

func (instance *Instance) callMethod(method *Function, class *Class, arguments []Object) Object {
	methodEnvironment := createMethodEnvironment(method, arguments)
	methodScope := &Scope{Self: instance, Class: class, Environment: methodEnvironment}

	return evaluator(method.Body, methodScope)
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
