package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateClass(node *ast.Class, scope *object.Scope) object.Object {
	class := &object.Class{
		Name:        node.Name,
		Scope:       scope,
		Environment: object.NewEnvironment(),
		Super:       nil,
	}

	// Create a new scope for this class
	classEnvironment := object.NewEnclosedEnvironment(scope.Environment)
	classScope := &object.Scope{Environment: classEnvironment, Self: class}

	Evaluate(node.Body, classScope)

	scope.Environment.Set(node.Name.Value, class)

	return class
}
