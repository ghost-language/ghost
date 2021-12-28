package interpreter

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/error"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func Evaluate(node ast.Node, env *object.Environment) (object.Object, bool) {
	switch node := node.(type) {
	case *ast.Program:
		return evaluateProgram(node, env)
	case *ast.Assign:
		return evaluateAssign(node, env)
	case *ast.Block:
		return evaluateBlock(node, env)
	case *ast.Boolean:
		return evaluateBoolean(node, env)
	case *ast.Call:
		return evaluateCall(node, env)
	case *ast.Expression:
		return Evaluate(node.Expression, env)
	case *ast.Function:
		return evaluateFunction(node, env)
	case *ast.Identifier:
		return evaluateIdentifier(node, env)
	case *ast.If:
		return evaluateIf(node, env)
	case *ast.Index:
		return evaluateIndex(node, env)
	case *ast.Infix:
		return evaluateInfix(node, env)
	case *ast.List:
		return evaluateList(node, env)
	case *ast.Map
		return evaluateMap(node, env)
	case *ast.Null:
		return evaluateNull(node, env)
	case *ast.Number:
		return evaluateNumber(node, env)
	case *ast.Prefix:
		return evaluatePrefix(node, env)
	case *ast.String:
		return evaluateString(node, env)
	case nil:
		return nil, false
	default:
		if node != nil {
			err := error.Error{
				Reason:  error.Runtime,
				Message: fmt.Sprintf("unknown runtime node: %T", node),
			}

			log.Error(err.Reason, err.Message)
		}

		return nil, false
	}
}

func toBooleanValue(input bool) *object.Boolean {
	if input {
		return value.TRUE
	}

	return value.FALSE
}

func isTruthy(value object.Object) bool {
	switch value := value.(type) {
	case *object.Null:
		return false
	case *object.Boolean:
		return value.Value
	case *object.String:
		return len(value.Value) > 0
	default:
		return true
	}
}
