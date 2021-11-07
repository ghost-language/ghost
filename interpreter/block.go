package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/object"
)

func evaluateBlock(node *ast.Block) (object.Object, bool) {
	log.Debug("evaluating block: %q", node.Statements)
	var result object.Object

	for _, statement := range node.Statements {
		result, _ = Evaluate(statement)
		log.Debug("block result: %q", result)
	}

	return result, true
}
