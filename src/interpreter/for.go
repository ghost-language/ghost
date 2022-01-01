package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateFor(node *ast.For, env *object.Environment) (object.Object, bool) {
	existingIdentifier, identifierExisted := env.Get(node.Identifier.Value)

	_, ok := Evaluate(node.Initializer, env)

	if !ok {
		return nil, false
	}

	loop := true

	defer func() {
		if identifierExisted {
			env.Set(node.Identifier.Value, existingIdentifier)
		} else {
			env.Delete(node.Identifier.Value)
		}
	}()

	for loop {
		condition, ok := Evaluate(node.Condition, env)

		if !ok {
			return nil, false
		}

		if isTruthy(condition) {
			_, ok = Evaluate(node.Block, env)

			if !ok {
				return nil, false
			}

			_, ok = Evaluate(node.Increment, env)

			if !ok {
				return nil, false
			}

			continue
		}

		loop = false
	}

	return nil, true
}
