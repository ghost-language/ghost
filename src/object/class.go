package object

import "ghostlang.org/x/ghost/ast"

const CLASS = "CLASS"

// Class objects consist of a body and an environment.
type Class struct {
	Name        *ast.Identifier
	Scope       *Scope
	Environment *Environment
	Super       *Class
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
	return nil, false
}
