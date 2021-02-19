package interpreter

import "ghostlang.org/x/ghost/ast"

// Evaluate parses the abstract syntax tree, evaluating each type of node and
// producing a result.
func Evaluate(expression ast.ExpressionNode) interface{} {
	switch node := expression.(type) {
	case *ast.Grouping:
		return Evaluate(node.Expression)
	case *ast.Literal:
		return node.Value
	}

	panic("Fatal error")
}
