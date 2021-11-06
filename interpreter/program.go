package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateProgram(node *ast.Program) (object.Object, bool) {
	var result object.Object
	var ok bool

	for _, statement := range node.Statements {
		result, ok = Evaluate(statement)

		if !ok {
			return nil, false
		}
	}

	return result, true
}
