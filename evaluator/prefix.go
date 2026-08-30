package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

func evaluatePrefix(node *ast.Prefix, scope *object.Scope) object.Object {
	right := Evaluate(node.Right, scope)

	if isError(right) {
		return right
	}

	switch node.Operator {
	case token.BANG:
		// Routes through object.IsFalse rather than hand-rolling truthiness
		// here (§13.11) - a hand-rolled version of this exact switch used to
		// live here, and had already drifted from object/boolean.go's
		// isTruthy: its default branch answered false for every type but
		// Boolean/Null, silently skipping the case String has in the real
		// rule (empty string is falsy), so !"" answered false instead of
		// true.
		return toBooleanValue(object.IsFalse(right))
	case token.MINUS:
		// Only works with number objects
		if right.Type() != object.NUMBER {
			return object.NewError(fault.Type, node.Token, "cannot negate %s", object.TypeName(right))
		}

		return right.(*object.Number).Neg()
	}

	return object.NewError(fault.Type, node.Token, "cannot use `%s` on %s", node.Operator, object.TypeName(right))
}
