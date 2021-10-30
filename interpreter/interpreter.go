package interpreter

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/error"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/object"
)

func Interpret(statements []ast.StatementNode) {
	for _, statement := range statements {
		result, ok := Evaluate(statement)

		if !ok {
			err := error.Error{
				Reason:  error.Runtime,
				Message: fmt.Sprintf("unknown runtime node: %T", statement),
			}

			log.Error(err.Reason, err.Message)
		} else {
			// temporarily log the returned object
			log.Info(fmt.Sprintf("== %s", result.String()))
		}
	}
}

func Evaluate(node ast.Node) (object.Object, bool) {
	switch node := node.(type) {
	case *ast.Expression:
		return Evaluate(node.Expression)
	case *ast.Number:
		return &object.Number{Value: node.Value}, true
	}

	return nil, false
}
