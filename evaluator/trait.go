package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateTrait(node *ast.Trait, scope *object.Scope) object.Object {
	trait := &object.Trait{
		Name:        node.Name,
		Scope:       scope,
		Environment: object.NewEnvironment(),
	}

	// Create a new scope for this trait
	trait.Environment = object.NewEnclosedEnvironment(scope.Environment)
	traitScope := &object.Scope{Environment: trait.Environment, Self: trait}

	Evaluate(node.Body, traitScope)

	scope.Environment.Set(node.Name.Value, trait)

	return trait
}
