package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateTrait(node *ast.Trait, scope *object.Scope) object.Object {
	// As with a class: the trait keeps this scope for its methods to read
	// through.
	scope.Environment.Capture()

	trait := &object.Trait{
		Name:        node.Name,
		Scope:       scope,
		Environment: object.NewEnvironment(),
	}

	// Create a new scope for this trait
	trait.Environment = object.NewEnclosedEnvironment(scope.Environment)
	traitScope := &object.Scope{Environment: trait.Environment, Self: trait}

	result := Evaluate(node.Body, traitScope)

	if isError(result) {
		return result
	}

	scope.Environment.Set(node.Name.Value, trait)

	return trait
}
