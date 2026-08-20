package evaluator

import (
	"ghostlang.org/x/ghost/ast"
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
		return newError("%d:%d:%s: runtime error: 'super' used outside of class context", node.Token.Line, node.Token.Column, node.Token.File)
	}

	if scope.Class == nil || scope.Class.Super == nil {
		return newError("%d:%d:%s: runtime error: class %s has no superclass", node.Token.Line, node.Token.Column, node.Token.File, instance.Class.Name.Value)
	}

	return &object.Super{Instance: instance, Class: scope.Class.Super}
}
