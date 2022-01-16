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

	// super
	if node.Super != nil {
		identifier, ok := scope.Environment.Get(node.Super.Value)

		if !ok {
			object.NewError("%d:%d: runtime error: identifier '%s' not found in '%s'", node.Super.Token.Line, node.Super.Token.Column, node.Super.Value, scope.Self.String())
		}

		super, ok := identifier.(*object.Class)

		if !ok {
			object.NewError("%d:%d: runtime error: referenced identifier in extends not a class, got=%T", node.Super.Token.Line, node.Super.Token.Column, super)
		}

		class.Super = super
	}

	// Create a new scope for this class
	classEnvironment := object.NewEnclosedEnvironment(scope.Environment)
	classScope := &object.Scope{Environment: classEnvironment, Self: class}

	Evaluate(node.Body, classScope)

	scope.Environment.Set(node.Name.Value, class)

	return class
}
