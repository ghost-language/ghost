package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateAssign(node *ast.Assign, env *object.Environment) (object.Object, bool) {
	value, _ := Evaluate(node.Value, env)

	if isError(value) {
		return value, false
	}

	if node.Name != nil {
		env.Set(node.Name.Value, value)
	}

	return nil, true
}
