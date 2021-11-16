package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/environment"
	"ghostlang.org/x/ghost/object"
)

func evaluateProgram(node *ast.Program, env *environment.Environment) (object.Object, bool) {
	var result object.Object

	for _, statement := range node.Statements {
		result, _ = Evaluate(statement, env)
	}

	return result, true
}
