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

// HasField reports whether this trait declares a field by that name (§13.18).
func (trait *Trait) HasField(name string) bool {
	return hasField(trait.Fields, name)
}

// DeclaredName names this trait for a diagnostic about its own body.
func (trait *Trait) DeclaredName() string {
	return trait.Name.Value
}
