package object

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

// Trait objects consist of a body and an environment.
type Trait struct {
	Name        *ast.Identifier
	Scope       *Scope
	Environment *Environment
	Fields      []Field
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
func (trait *Trait) Method(method string, tok token.Token, args []Object) (Object, bool) {
	return nil, false
}

// SetField records a field declaration on the trait.
func (trait *Trait) SetField(name string, value ast.ExpressionNode) {
	trait.Fields = setField(trait.Fields, name, value)
}
