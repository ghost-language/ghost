package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

func evaluateInfix(node *ast.Infix, scope *object.Scope) object.Object {
	left := Evaluate(node.Left, scope)

	if isError(left) {
		return left
	}

	right := Evaluate(node.Right, scope)

	if isError(right) {
		return right
	}

	if left == nil || right == nil {
		return newError("%d:%d:%s: runtime error: operand of %s is null", node.Token.Line, node.Token.Column, node.Token.File, node.Operator)
	}

	if node.Operator == token.EQUALEQUAL || node.Operator == token.BANGEQUAL {
		if result, ok := evaluateEquality(node, left, right); ok {
			return result
		}
	}

	switch {
	case left.Type() == object.BOOLEAN && right.Type() == object.BOOLEAN:
		return evaluateBooleanInfix(node, left, right)
	case left.Type() == object.NUMBER && right.Type() == object.NUMBER:
		return evaluateNumberInfix(node, left, right)
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return evaluateStringInfix(node, left, right)
	case left.Type() != right.Type():
		return newError("%d:%d:%s: runtime error: type mismatch: %s %s %s", node.Token.Line, node.Token.Column, node.Token.File, left.Type(), node.Operator, right.Type())
	}

	return newError("%d:%d:%s: runtime error: unknown operator: %s %s %s", node.Token.Line, node.Token.Column, node.Token.File, left.Type(), node.Operator, right.Type())
}

// evaluateEquality handles the comparisons that are not type-specific: a
// comparison against null, which is how code tests whether a value is set, and
// identity between two instances. Operands are otherwise required to share a
// type, so both of these would fall through to a type mismatch error.
//
// It reports false when the operands are not its business, leaving the
// type-specific comparisons below to run.
func evaluateEquality(node *ast.Infix, left object.Object, right object.Object) (object.Object, bool) {
	var equal bool

	switch {
	case left.Type() == object.NULL || right.Type() == object.NULL:
		equal = left.Type() == object.NULL && right.Type() == object.NULL
	case left.Type() == object.INSTANCE && right.Type() == object.INSTANCE:
		// Instances compare by identity: two instances of the same class with
		// equal fields are still two different objects.
		equal = left == right
	default:
		return nil, false
	}

	if node.Operator == token.BANGEQUAL {
		return toBooleanValue(!equal), true
	}

	return toBooleanValue(equal), true
}
