package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
	"ghostlang.org/x/ghost/value"
)

func evaluatePrefix(node *ast.Prefix, scope *object.Scope) object.Object {
	right := Evaluate(node.Right, scope)

	if isError(right) {
		return right
	}

	switch node.Operator {
	case token.BANG:
		// Compare by value rather than by pointer identity. Not every boolean
		// reaching this point is one of the value.TRUE/value.FALSE singletons:
		// string comparisons and library functions such as string.startsWith()
		// and math.isNegative() build fresh boolean objects, and an identity
		// check silently fell through to the default branch for those, making
		// !(expression) yield false regardless of the operand.
		switch right := right.(type) {
		case *object.Boolean:
			return toBooleanValue(!right.Value)
		case *object.Null:
			return value.TRUE
		default:
			return value.FALSE
		}
	case token.MINUS:
		// Only works with number objects
		if right.Type() != object.NUMBER {
			return newError("%d:%d:%s: runtime error: unknown operator: -%s", node.Token.Line, node.Token.Column, node.Token.File, right.Type())
		}

		return right.(*object.Number).Neg()
	}

	return newError("%d:%d:%s: runtime error: unknown operator: %s%s", node.Token.Line, node.Token.Column, node.Token.File, node.Operator, right.Type())
}
