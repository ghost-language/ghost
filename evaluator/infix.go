package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
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
		return object.NewError(fault.Type, node.Token, "cannot use `%s` when one side is null", node.Operator)
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
	case left.Type() == object.DATE && right.Type() == object.DATE:
		return evaluateDateInfix(node, left, right)
	case isListArithmetic(left, right):
		return evaluateListInfix(node, left, right)
	case left.Type() != right.Type():
		return object.NewError(fault.Type, node.Token, "cannot use `%s` between %s and %s", node.Operator, object.TypeName(left), object.TypeName(right))
	}

	return operatorError(node.Token, node.Operator, left, right)
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
	case left.Type() == object.LIST && right.Type() == object.LIST:
		// Lists compare by contents, to any depth. Two lists written out the
		// same way are equal, which is the only reading that makes `==` useful
		// on a value built rather than passed around.
		equal = listsEqual(left.(*object.List), right.(*object.List))
	default:
		return nil, false
	}

	if node.Operator == token.BANGEQUAL {
		return toBooleanValue(!equal), true
	}

	return toBooleanValue(equal), true
}

// isListArithmetic reports whether an operation is elementwise list arithmetic.
// A list against a list or against a number is; a list against anything else
// stays a type mismatch, which says more than a broadcasting failure would.
func isListArithmetic(left object.Object, right object.Object) bool {
	if left.Type() == object.LIST {
		return right.Type() == object.LIST || right.Type() == object.NUMBER
	}

	return left.Type() == object.NUMBER && right.Type() == object.LIST
}
