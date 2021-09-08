package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateGet(node *ast.Get, env *object.Environment) (object.Object, bool) {
	value, ok := Evaluate(node.Expression, env)

	if !ok {
		return nil, ok
	}

	if accessor, ok := value.(object.PropertyAccessor); ok {
		return accessor.Get(node.Name), true
	}

	return nil, false
}