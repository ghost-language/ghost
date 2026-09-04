package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

func evaluateBoolean(node *ast.Boolean, scope *object.Scope) object.Object {
	return toBooleanValue(node.Value)
}

func evaluateBooleanInfix(node *ast.Infix, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.Boolean).Value
	rightValue := right.(*object.Boolean).Value

	// `and` and `or` are absent deliberately: evaluateInfix routes them to
	// evaluateLogicalInfix before this function is reached, because they have
	// to decide whether to evaluate the right operand at all and everything
	// here is handed both operands already evaluated.
	switch node.Operator {
	case token.EQUALEQUAL:
		return toBooleanValue(leftValue == rightValue)
	case token.BANGEQUAL:
		return toBooleanValue(leftValue != rightValue)
	}

	return object.NewError(fault.Type, node.Token, "cannot use `%s` between two booleans", node.Operator)
}

// evaluateLogicalInfix evaluates `and` and `or`, which short-circuit: the
// right operand is evaluated only when the left one leaves the answer open
// (§13.21, §14 decision 11). `false and x` is false and `true or x` is true
// whatever x turns out to be, so x is never reached - which is what lets a
// guard like `x == null or x.field` answer before the dereference it exists
// to prevent.
//
// Both operands are still booleans, exactly as before. What changes is only
// which operands there are: one that is reached and is not a boolean is a
// Type fault as it always was, and one that is never reached is never
// checked.
func evaluateLogicalInfix(node *ast.Infix, left object.Object, scope *object.Scope) object.Object {
	condition, ok := left.(*object.Boolean)

	if !ok {
		return logicalOperandError(node.Token, node.Operator, left, "left")
	}

	// The left operand settles it on its own. Returning here without touching
	// node.Right is the whole of the short-circuit.
	if node.Operator == token.AND && !condition.Value {
		return toBooleanValue(false)
	}

	if node.Operator == token.OR && condition.Value {
		return toBooleanValue(true)
	}

	right := Evaluate(node.Right, scope)

	if isError(right) {
		return right
	}

	result, ok := right.(*object.Boolean)

	if !ok {
		return logicalOperandError(node.Token, node.Operator, right, "right")
	}

	// The left operand did not settle it, so the answer is whatever the right
	// one says: `true and x` is x, and `false or x` is x.
	return toBooleanValue(result.Value)
}
