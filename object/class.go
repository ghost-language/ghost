package object

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

// Field is a class or trait field declaration. The initializer is kept
// unevaluated so that every instance gets its own value rather than sharing one
// object stored on the class.
type Field struct {
	Name  string
	Value ast.ExpressionNode
}

// FieldDeclarer is implemented by the declarations that can hold fields, so the
// evaluator can record a field without caring whether it is on a class or a
// trait.
type FieldDeclarer interface {
	SetField(name string, value ast.ExpressionNode)
	HasField(name string) bool
	DeclarationScope() *Scope
}

// Class objects consist of a body and an environment.
type Class struct {
	Name        *ast.Identifier
	Scope       *Scope
	Environment *Environment
	Super       *Class
	Traits      []*Trait
	Fields      []Field
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
func (class *Class) Method(method string, tok token.Token, args []Object) (Object, bool) {
	return nil, false
}

// SetField records a field declaration, replacing any earlier declaration of
// the same name so a redeclaration in the same body wins.
func (class *Class) SetField(name string, value ast.ExpressionNode) {
	class.Fields = setField(class.Fields, name, value)
}

// HasField reports whether this class declares a field of its own by that
// name, which is what lets a method declaration notice it is colliding with
// one (§13.18).
func (class *Class) HasField(name string) bool {
	return hasField(class.Fields, name)
}

// DeclarationScope is the scope the class was declared in, which is what its
// methods close over. A method is a member reached through `this`, not a
// lexical binding, so its body resolves past the class's own member table -
// see §14 decision 12.
func (class *Class) DeclarationScope() *Scope {
	return class.Scope
}

// Ancestors returns the class chain from the root superclass down to this
// class, which is the order fields and constructors must be applied in.
func (class *Class) Ancestors() []*Class {
	var chain []*Class

	for current := class; current != nil; current = current.Super {
		chain = append([]*Class{current}, chain...)
	}

	return chain
}

func hasField(fields []Field, name string) bool {
	for _, field := range fields {
		if field.Name == name {
			return true
		}
	}

	return false
}

func setField(fields []Field, name string, value ast.ExpressionNode) []Field {
	for index, field := range fields {
		if field.Name == name {
			fields[index].Value = value

			return fields
		}
	}

	return append(fields, Field{Name: name, Value: value})
}
