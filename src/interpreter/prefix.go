package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluatePrefix(node *ast.Prefix, env *object.Environment) object.Object {
	right := Evaluate(node.Right, env)

	if isError(right) {
		return right
	}

	switch node.Operator {
	case "!":
		switch right {
		case value.TRUE:
			return value.FALSE
		case value.FALSE:
			return value.TRUE
		case value.NULL:
			return value.TRUE
		default:
			return value.FALSE
		}
	case "-":
		// Only works with number objects
		if right.Type() != object.NUMBER {
			return newError("unknown operator: -%s", right.Type())
		}

		numberValue := right.(*object.Number).Value.Neg()

		return &object.Number{Value: numberValue}
	}

	return newError("unknown operator: %s%s", node.Operator, right.Type())
}
