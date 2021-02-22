package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

// Evaluate parses the abstract syntax tree, evaluating each type of node and
// producing a result.
func Evaluate(expression ast.ExpressionNode) object.Object {
	switch node := expression.(type) {
	case *ast.Binary:
		return evaluateBinary(node)
	case *ast.Grouping:
		return Evaluate(node.Expression)
	case *ast.Number:
		return &object.Number{Value: node.Value}
	case *ast.String:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		return &object.Boolean{Value: node.Value}
	case *ast.Null:
		return &object.Null{}
	}

	panic("Fatal error")
}
