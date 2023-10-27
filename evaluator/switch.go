package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateSwitch(node *ast.Switch, scope *object.Scope) object.Object {
	// Get the value
	obj := Evaluate(node.Value, scope)

	for _, option := range node.Cases {
		// Skip default case to handle last if needed
		if option.Default {
			continue
		}

		for _, val := range option.Value {
			out := Evaluate(val, scope)

			if obj.Type() == out.Type() && (obj.String() == out.String()) {
				// evaluate the block and return the value
				out := evaluateBlock(option.Body, scope)

				return out
			}
		}
	}

	// Handle default case
	for _, option := range node.Cases {
		if option.Default {
			out := evaluateBlock(option.Body, scope)

			return out
		}
	}

	return nil
}
