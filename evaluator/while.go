package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateWhile(node *ast.While, scope *object.Scope) object.Object {
	for {
		condition := Evaluate(node.Condition, scope)

		if isError(condition) {
			return condition
		}

		if object.IsTrue(condition) {
			// A scope per iteration, so a closure made in the body captures
			// that iteration's bindings rather than one shared set. When no
			// closure took it, the environment comes straight back for the
			// next iteration to reuse.
			iteration := enclose(scope)
			evaluated := Evaluate(node.Consequence, iteration)

			release(scope, iteration)

			if isTerminator(evaluated) {
				switch val := evaluated.(type) {
				case *object.Error:
					return val
				case *object.Return:
					return val
				case *object.Continue:
					//
				case *object.Break:
					return nil
				}
			}
		} else {
			break
		}
	}

	return nil
}
