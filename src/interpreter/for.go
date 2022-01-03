package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateFor(node *ast.For, env *object.Environment) object.Object {
	existingIdentifier, identifierExisted := env.Get(node.Identifier.Value)

	defer func() {
		if identifierExisted {
			env.Set(node.Identifier.Value, existingIdentifier)
		} else {
			env.Delete(node.Identifier.Value)
		}
	}()

	initializer := Evaluate(node.Initializer, env)

	if isError(initializer) {
		return initializer
	}

	loop := true

	for loop {
		condition := Evaluate(node.Condition, env)

		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			err := Evaluate(node.Block, env)

			if isError(err) {
				return err
			}

			err = Evaluate(node.Increment, env)

			if isError(err) {
				return err
			}

			continue
		}

		loop = false
	}

	return nil
}
