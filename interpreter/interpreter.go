package interpreter

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/error"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func Interpret(statements []ast.StatementNode) {
	for _, statement := range statements {
		_, ok := Evaluate(statement)

		if ok {
			// temporarily log the returned object
			// log.Info(fmt.Sprintf("== %s", result.String()))
		}
	}
}

func Evaluate(node ast.Node) (object.Object, bool) {
	switch node := node.(type) {
	case *ast.Boolean:
		return evaluateBoolean(node)
	case *ast.Call:
		return evaluateCall(node)
	case *ast.Expression:
		return Evaluate(node.Expression)
	case *ast.Identifier:
		return evaluateIdentifier(node)
	case *ast.Infix:
		return evaluateInfix(node)
	case *ast.Null:
		return evaluateNull(node)
	case *ast.Number:
		return evaluateNumber(node)
	case *ast.Prefix:
		return evaluatePrefix(node)
	case *ast.String:
		return evaluateString(node)
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
