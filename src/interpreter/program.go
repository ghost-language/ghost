package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateProgram(node *ast.Program, env *object.Environment) (object.Object, bool) {
	var result object.Object

	for _, statement := range node.Statements {
		result, _ = Evaluate(statement, env)
	}

	return result, true
}
