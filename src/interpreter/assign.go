package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/environment"
	"ghostlang.org/x/ghost/object"
)

func evaluateAssign(node *ast.Assign, env *environment.Environment) (object.Object, bool) {
	value, ok := Evaluate(node.Value, env)

	if !ok {
		return nil, false
	}

	if node.Name != nil {
		env.Set(node.Name.Value, value)
	}

	return nil, true
}
