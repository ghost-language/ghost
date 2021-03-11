package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/environment"
	"ghostlang.org/x/ghost/object"
)

func evaluateAssign(node *ast.Assign, env *environment.Environment) (object.Object, bool) {
	val, _ := Evaluate(node.Value, env)

	env.Set(node.Name.Lexeme, val)

	return nil, true
}
