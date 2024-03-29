package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateList(node *ast.List, scope *object.Scope) object.Object {
	elements := evaluateExpressions(node.Elements, scope)

	if len(elements) == 1 && isError(elements[0]) {
		return elements[0]
	}

	return &object.List{Elements: elements}
}
