package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/environment"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluatePrefix(node *ast.Prefix, env *environment.Environment) (object.Object, bool) {
	right, ok := Evaluate(node.Right, env)

	if !ok {
		return nil, false
	}

	switch node.Operator {
	case "!":
		switch right {
		case value.TRUE:
			return value.FALSE, true
		case value.FALSE:
			return value.TRUE, true
		case value.NULL:
			return value.TRUE, true
		default:
			return value.FALSE, true
		}
	case "-":
		// Only works with number objects
		if right.Type() != object.NUMBER {
			// TODO: error message
			return nil, false
		}

		numberValue := right.(*object.Number).Value.Neg()

		return &object.Number{Value: numberValue}, true
	default:
		return nil, false
	}
}
