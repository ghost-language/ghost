package interpreter

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluatePrint(node *ast.Print) (object.Object, bool) {
	result, ok := Evaluate(node.Expression)

	if !ok {
		return result, ok
	}

	fmt.Println(result.String())

	return value.NULL, true
}
