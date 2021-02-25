package interpreter

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/ghost"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

// Interpret ...
func Interpret(statements []ast.StatementNode) {
	for _, statement := range statements {
		result, ok := Evaluate(statement)

		if !ok {
			ghost.RuntimeError(result.String())
		}
	}
}

// Evaluate parses the abstract syntax tree, evaluating each type of node and
// producing a result.
func Evaluate(node ast.Node) (object.Object, bool) {
	switch node := node.(type) {
	case *ast.Binary:
		return evaluateBinary(node)
	case *ast.Boolean:
		if node.Value {
			return value.TRUE, true
		}

		return value.FALSE, true
	case *ast.Expression:
		result, ok := Evaluate(node.Expression)

		if !ok {
			return result, ok
		}

		return value.NULL, ok
	case *ast.Grouping:
		return Evaluate(node.Expression)
	case *ast.Null:
		return value.NULL, true
	case *ast.Number:
		return &object.Number{Value: node.Value}, true
	case *ast.Print:
		return evaluatePrint(node)
	case *ast.String:
		return &object.String{Value: node.Value}, true
	case *ast.Ternary:
		return evaluateTernary(node)
	case *ast.Unary:
		return evaluateUnary(node)
	}

	return &object.Error{Message: fmt.Sprintf("unrecognized node: %v", node)}, false
}
