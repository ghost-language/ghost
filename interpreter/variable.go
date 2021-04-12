package interpreter

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/standard"
)

func evaluateIdentifier(node *ast.Identifier, env *object.Environment) (object.Object, bool) {
	val, err := env.Get(node.Name)

	if err != nil {
		standard, success := standard.StandardFunctions[node.Name.Lexeme]

		if success != true {
			return &object.Error{Message: fmt.Sprintf("unknown identifier: %s", node.Name.Lexeme)}, false
		}

		return standard, true
	}

	return val, true
}