package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

// Evaluate parses the abstract syntax tree, evaluating each type of node and
// producing a result.
func Evaluate(expression ast.ExpressionNode) object.Object {
	switch node := expression.(type) {
	case *ast.Binary:
		return evaluateBinary(node)
	case *ast.Boolean:
		if node.Value {
			return value.TRUE
		}

		return value.FALSE
	case *ast.Grouping:
		return Evaluate(node.Expression)
	case *ast.Null:
		return value.NULL
	case *ast.Number:
		return &object.Number{Value: node.Value}
	case *ast.String:
		return &object.String{Value: node.Value}
	case *ast.Unary:
		return evaluateUnary(node)
	}

	panic("Fatal error")
}
