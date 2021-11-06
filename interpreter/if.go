package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/object"
)

func evaluateIf(node *ast.If) (object.Object, bool) {
	log.Debug("evaluating if statement")
	condition, ok := Evaluate(node.Condition)

	if !ok {
		return nil, false
	}

	var result object.Object

	if isTruthy(condition) {
		log.Debug("condition true")
		result, _ = Evaluate(node.Consequence)
		log.Debug("result: %s", result)
	} else if node.Alternative != nil {
		log.Debug("condition false")
		result, _ = Evaluate(node.Consequence)
	}

	return result, true
}
