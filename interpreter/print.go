package interpreter

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluatePrint(node *ast.Print, env *object.Environment) (object.Object, bool) {
	result, ok := Evaluate(node.Expression, env)

	if !ok {
		return result, ok
	}

	fmt.Fprintln(env.GetWriter(), result.String())

	return value.NULL, true
}
