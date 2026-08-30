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
				return evaluateBlock(option.Body, scope)
			}
		}
	}

	// Handle default case
	for _, option := range node.Cases {
		if option.Default {
			return evaluateBlock(option.Body, scope)
		}
	}

	return nil
}
