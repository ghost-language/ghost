package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
)

func evaluateClass(node *ast.Class, scope *object.Scope) object.Object {
	// The class keeps this scope, and its methods read through it, so it
	// outlives the block it was declared in.
	scope.Environment.Capture()

	class := &object.Class{
		Name:  node.Name,
		Scope: scope,
		Super: nil,
	}

	// super
	if node.Super != nil {
		identifier, ok := scope.Environment.Get(node.Super.Value)

		if !ok {
			return object.NewError(fault.Name, node.Super.Token, "`%s` is not defined", node.Super.Value).
				WithHelp("a class has to be declared before the class that extends it")
		}

		super, ok := identifier.(*object.Class)

		if !ok {
			return object.NewError(fault.Type, node.Super.Token, "cannot extend `%s`, which is a %s, not a class", node.Super.Value, object.TypeName(identifier))
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
