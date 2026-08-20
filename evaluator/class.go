package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateClass(node *ast.Class, scope *object.Scope) object.Object {
	class := &object.Class{
		Name:  node.Name,
		Scope: scope,
		Super: nil,
	}

	// super
	if node.Super != nil {
		identifier, ok := scope.Environment.Get(node.Super.Value)

		if !ok {
			return object.NewError("%d:%d:%s: runtime error: unknown identifier: %s", node.Super.Token.Line, node.Super.Token.Column, node.Super.Token.File, node.Super.Value)
		}

		super, ok := identifier.(*object.Class)

		if !ok {
			return object.NewError("%d:%d:%s: runtime error: referenced identifier in extends not a class, got=%s", node.Super.Token.Line, node.Super.Token.Column, node.Super.Token.File, identifier.Type())
		}

		class.Super = super
	}

	// The class environment holds the class's members and doubles as the scope
	// its body is evaluated in, so a method can call a sibling method by bare
	// name. Enclosing it in the defining scope keeps outer bindings reachable
	// from method bodies.
	class.Environment = object.NewEnclosedEnvironment(scope.Environment)
	classScope := &object.Scope{Environment: class.Environment, Self: class, Class: class}

	result := Evaluate(node.Body, classScope)

	if isError(result) {
		return result
	}

	scope.Environment.Set(node.Name.Value, class)

	return class
}
