package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateSwitch(node *ast.Switch, scope *object.Scope) object.Object {
	obj := Evaluate(node.Value, scope)

	if isError(obj) {
		return obj
	}

	for _, option := range node.Cases {
		// Skip default case to handle last if needed
		if option.Default {
			continue
		}

		for _, val := range option.Value {
			out := Evaluate(val, scope)

			if isError(out) {
				return out
			}

			// The same rule `==` uses everywhere else in the language (§13.2):
			// content equality for a list/map/duration/date, identity for
			// everything else. Comparing by String() used to accept any two
			// functions or classes as a match - every Function stringifies to
			// the literal "function" and every Class to "class" - regardless
			// of which one either side actually was.
			if object.ValuesEqual(obj, out) {
				return evaluateCase(option.Body, scope)
			}
		}
	}

	// Handle default case
	for _, option := range node.Cases {
		if option.Default {
			return evaluateCase(option.Body, scope)
		}
	}

	return nil
}

// evaluateCase runs a matched case's body in a scope of its own.
func evaluateCase(body *ast.Block, scope *object.Scope) object.Object {
	block := enclose(scope)
	result := evaluateBlock(body, block)

	release(scope, block)

	return result
}
