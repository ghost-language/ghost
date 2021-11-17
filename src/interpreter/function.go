package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateFunction(node *ast.Function, env *object.Environment) (object.Object, bool) {
	function := &object.Function{
		Parameters:  node.Parameters,
		Defaults:    node.Defaults,
		Body:        node.Body,
		Environment: env,
	}

	if node.Name != nil {
		env.Set(node.Name.Value, function)
	}

	return function, true
}
