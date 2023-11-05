package object

import "ghostlang.org/x/ghost/ast"

const TRAIT = "TRAIT"

// Trait objects consist of a body and an environment.
type Trait struct {
	Name        *ast.Identifier
	Scope       *Scope
	Environment *Environment
}

// String represents the class object's value as a string.
func (trait *Trait) String() string {
	return "trait"
}

// Type returns the trait object type.
func (trait *Trait) Type() Type {
	return TRAIT
}

// Method defines the set of methods available on trait objects.
func (trait *Trait) Method(method string, args []Object) (Object, bool) {
	return nil, false
}
