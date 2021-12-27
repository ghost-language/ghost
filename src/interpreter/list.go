package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateList(node *ast.List, env *object.Environment) (object.Object, bool) {
	elements, ok := evaluateExpressions(node.Elements, env)

	if !ok {
		return nil, false
	}

	return &object.List{Elements: elements}, true
}
