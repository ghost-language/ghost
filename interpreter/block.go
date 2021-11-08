package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateBlock(node *ast.Block) (object.Object, bool) {
	var result object.Object

	for _, statement := range node.Statements {
		result, _ = Evaluate(statement)
	}

	return result, true
}
