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

// evaluateEquality handles every `==`/`!=` comparison except the four types
// that read their own operator instead: Boolean, Number, String, and Date
// each have a reason of their own to keep their dedicated infix evaluator
// (evaluateNumberInfix promotes int/float the way `+` does; evaluateDateInfix
// compares the instant rather than every field, so a Date attached to two
// different time zones can still be equal) - a same-typed pair of those
// reports false here and falls through to that evaluator instead.
//
// Everything else goes through object.ValuesEqual: a comparison against
// null, which is how code tests whether a value is set, and a same-type
// comparison for every other type - list/map/duration by content, everything
// else (instances, functions, classes, ...) by identity, matching the
// identity comparison this function used to give Instance alone (§13.2). A
// comparison between two different, non-null types is still nobody's
// business here - it reports false too, so `evaluateInfix`'s type-mismatch
// error fires instead of silently answering false, the same as it always has
// for `1 == "a"`.
func evaluateEquality(node *ast.Infix, left object.Object, right object.Object) (object.Object, bool) {
	sameType := left.Type() == right.Type()
	involvesNull := left.Type() == object.NULL || right.Type() == object.NULL

	if !sameType && !involvesNull {
		return nil, false
	}

	if sameType {
		switch left.Type() {
		case object.BOOLEAN, object.NUMBER, object.STRING, object.DATE:
			return nil, false
		}
	}

	equal := object.ValuesEqual(left, right)

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
