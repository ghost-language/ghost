package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateClass(node *ast.Class, env *object.Environment) object.Object {
	class := &object.Class{
		Name:        node.Name,
		Environment: object.NewEnvironment(),
	}

	env.Set(node.Name.Value, class)

	return class
}
