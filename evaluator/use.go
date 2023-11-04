package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateUse(node *ast.Use, scope *object.Scope) object.Object {
	// check that the scope is a class
	class, ok := scope.Self.(*object.Class)

	if !ok {
		return object.NewError("%d:%d:%s: runtime error: use statement can only be used in a class", node.Token.Line, node.Token.Column, node.Token.File)
	}

	var traits []*object.Trait

	for _, trait := range node.Traits {
		if !scope.Environment.Has(trait.Value) {
			return object.NewError("%d:%d:%s: runtime error: trait '%s' is not defined", trait.Token.Line, trait.Token.Column, trait.Token.File, trait.Value)
		}

		identifier, _ := scope.Environment.Get(trait.Value)

		t, ok := identifier.(*object.Trait)

		if !ok {
			return object.NewError("%d:%d:%s: runtime error: referenced identifier in use not a trait, got=%T", trait.Token.Line, trait.Token.Column, trait.Token.File, trait)
		}

		traits = append(traits, t)
	}

	class.Traits = traits

	return nil
}
