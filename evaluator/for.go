package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateFor(node *ast.For, scope *object.Scope) object.Object {
	name := node.Identifier.Value

	// The control variable is carried between iterations here rather than in a
	// scope wrapping the loop, so the loop costs one level of environment
	// rather than two — every read of an outer variable from the body walks
	// that chain, and a loop nested four deep pays for each level.
	var carried object.Object = value.NULL
	started := false

	for {
		// Each iteration runs in a scope of its own, so a closure made in the
		// body captures the value that iteration saw rather than the one the
		// loop finished on (§13.14). The control variable is declared in it
		// directly, which both scopes it to the loop (§8.3) and stops the
		// initializer's assignment from walking outward (§13.13) and
		// overwriting a variable of the same name outside.
		iteration := enclose(scope)
		iteration.Environment.Set(name, carried)

		// The initializer runs once, and the increment runs at the top of
		// every iteration after it rather than at the foot of the previous
		// one. That ordering is what keeps a closure made in the body seeing
		// the value its own iteration ran with: incrementing in the iteration
		// the closure captured would move the value out from under it.
		if !started {
			if result := Evaluate(node.Initializer, iteration); isError(result) {
				return result
			}

			started = true
		} else if result := Evaluate(node.Increment, iteration); isError(result) {
			return result
		}

		condition := Evaluate(node.Condition, iteration)

		if isError(condition) {
			return condition
		}

		if !object.IsTrue(condition) {
			release(scope, iteration)

			return value.NULL
		}

		evaluated := Evaluate(node.Block, iteration)

		if current, ok := iteration.Environment.GetLocal(name); ok {
			carried = current
		}

		release(scope, iteration)

		if isTerminator(evaluated) {
			switch result := evaluated.(type) {
			case *object.Error:
				return result
			case *object.Return:
				return result
			case *object.Continue:
				//
			case *object.Break:
				return nil
			}
		}
	}
}
