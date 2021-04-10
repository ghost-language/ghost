package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/environment"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateBoolean(node *ast.Boolean, env *environment.Environment) (object.Object, bool) {
	if node.Value {
		return value.TRUE, true
	}

	return value.FALSE, true
}