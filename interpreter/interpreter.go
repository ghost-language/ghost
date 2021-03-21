package interpreter

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/environment"
	"ghostlang.org/x/ghost/errors"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/standard"
	"ghostlang.org/x/ghost/value"
)

// Interpret ...
func Interpret(statements []ast.StatementNode, env *environment.Environment) {
	for _, statement := range statements {
		result, ok := Evaluate(statement, env)

		if !ok {
			errors.RuntimeError(result.String())
		}
	}
}

// Evaluate parses the abstract syntax tree, evaluating each type of node and
// producing a result.
func Evaluate(node ast.Node, env *environment.Environment) (object.Object, bool) {
	switch node := node.(type) {
	case *ast.Assign:
		return evaluateAssign(node, env)
	case *ast.Binary:
		return evaluateBinary(node, env)
	case *ast.Block:
		return evaluateBlock(node, env)
	case *ast.Boolean:
		if node.Value {
			return value.TRUE, true
		}

		return value.FALSE, true
	case *ast.Call:
		return evaluateCall(node, env)
	case *ast.Declaration:
		return evaluateDeclaration(node, env)
	case *ast.Expression:
		result, ok := Evaluate(node.Expression, env)

		if !ok {
			return result, ok
		}

		return value.NULL, ok
	case *ast.Grouping:
		return Evaluate(node.Expression, env)
	case *ast.If:
		return evaluateIf(node, env)
	case *ast.Logical:
		return evaluateLogical(node, env)
	case *ast.Null:
		return value.NULL, true
	case *ast.Number:
		return &object.Number{Value: node.Value}, true
	case *ast.Print:
		return evaluatePrint(node, env)
	case *ast.String:
		return &object.String{Value: node.Value}, true
	case *ast.Ternary:
		return evaluateTernary(node, env)
	case *ast.Unary:
		return evaluateUnary(node, env)
	case *ast.Variable:
		val, err := env.Get(node.Name)

		if err != nil {
			standard, success := standard.StandardFunctions[node.Name.Lexeme]

			if success != true {
				return &object.Error{Message: fmt.Sprintf("unknown identifier: %v", node)}, false
			}

			return standard, true
		}

		return val, true
	case *ast.While:
		return evaluateWhile(node, env)
	}

	return &object.Error{Message: fmt.Sprintf("unrecognized node: %v", node)}, false
}
