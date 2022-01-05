package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateProgram(node *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range node.Statements {
		result = Evaluate(statement, env)

		switch statement := result.(type) {
		case *object.Error:
			return statement
		case *object.Return:
			return statement.Value
		}
	}

	return result
}