package object

import (
	"ghostlang.org/x/ghost/ast"
)

const CLASS = "CLASS"

// Class objects consist of a body and an environment.
type Class struct {
	Name        *ast.Identifier
	Scope       *Scope
	Environment *Environment
	Super       *Class
	Traits      []*Trait
}

// String represents the class object's value as a string.
func (class *Class) String() string {
	return "class"
}

// Type returns the class object type.
func (class *Class) Type() Type {
	return CLASS
}

// Method defines the set of methods available on class objects.
func (class *Class) Method(method string, args []Object) (Object, bool) {
	switch method {
	case "new":
		instance := &Instance{Class: class, Environment: NewEnclosedEnvironment(class.Environment)}

		// Check for constructor in class and superclass chain
		currentClass := class
		for currentClass != nil {
			if constructor, ok := currentClass.Environment.Get("constructor"); ok {
				result := instance.callMethod(constructor.(*Function), args, class.Name.Token)

				if result != nil && result.Type() == ERROR {
					return result, false
				}

				break
			}
			currentClass = currentClass.Super
		}

		return instance, true
	}

	return nil, false
}
