package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

// enclose returns a scope that runs in a fresh environment chained to this
// one, leaving everything else about the scope — the receiver, the declaring
// class, the call depth — exactly as it was.
//
// This is how a block gets a scope of its own (§13.15, §8.3). It is applied by
// the statements that own a block (`if`, `while`, `for`, `for ... in`,
// `switch`) rather than by evaluateBlock itself, because two of evaluateBlock's
// callers must not get one: a class or trait body is evaluated directly in the
// environment that collects its members, and a function or method body already
// runs in the frame createFunctionEnvironment built for it, which a second
// environment would only cost an allocation to wrap.
func enclose(scope *object.Scope) *object.Scope {
	return scope.Enclose()
}

// release hands a finished block scope's environment back to its parent for
// the next block to reuse. It is what keeps block scoping from allocating an
// environment per iteration of a hot loop; an environment a closure captured
// is kept rather than reused, which Environment.Release decides.
//
// Skipping it is only ever a missed reuse, never a correctness problem, so
// the paths that leave a block early — an error, a `return` — do not have to
// take care to call it.
func release(scope *object.Scope, block *object.Scope) {
	scope.Release(block)
}

func evaluateBlock(node *ast.Block, scope *object.Scope) object.Object {
	var result object.Object

	for _, statement := range node.Statements {
		result = Evaluate(statement, scope)

		if result != nil {
			resultType := result.Type()

			if resultType == object.ERROR || resultType == object.RETURN || resultType == object.CONTINUE || resultType == object.BREAK {
				return result
			}
		}
	}

	return result
}
