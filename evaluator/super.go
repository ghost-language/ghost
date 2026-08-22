package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
)

// evaluateSuper resolves `super` to a receiver bound to the current instance
// but starting member lookup at the superclass of the class that declared the
// running method. Anchoring on the declaring class rather than the instance's
// class is what keeps a `super` call in an inherited method from resolving back
// to itself.
func evaluateSuper(node *ast.Super, scope *object.Scope) object.Object {
	instance, ok := scope.Self.(*object.Instance)

	if !ok {
		return object.NewError(fault.Name, node.Token, "`super` can only be used inside a class")
	}

	if scope.Class == nil || scope.Class.Super == nil {
		return object.NewError(fault.Name, node.Token, "class `%s` has no superclass", instance.Class.Name.Value).
			WithHelp("declare one with `class %s extends Parent`", instance.Class.Name.Value)
	}

	return &object.Super{Instance: instance, Class: scope.Class.Super}
}
