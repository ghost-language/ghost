package interpreter

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/environment"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluatePrint(node *ast.Print, env *environment.Environment) (object.Object, bool) {
	result, ok := Evaluate(node.Expression, env)

	if !ok {
		return result, ok
	}

	fmt.Println(result.String())

	return value.NULL, true
}
