package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateFor(node *ast.For, scope *object.Scope) object.Object {
	existingIdentifier, identifierExisted := scope.Environment.Get(node.Identifier.Value)

	defer func() {
		if identifierExisted {
			scope.Environment.Set(node.Identifier.Value, existingIdentifier)
		} else {
			scope.Environment.Delete(node.Identifier.Value)
		}
	}()

	initializer := Evaluate(node.Initializer, scope)

	if isError(initializer) {
		return initializer
	}

	loop := true

	for loop {
		condition := Evaluate(node.Condition, scope)

		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			err := Evaluate(node.Block, scope)

			if isError(err) {
				return err
			}

			err = Evaluate(node.Increment, scope)

			if isError(err) {
				return err
			}

			continue
		}

		loop = false
	}

	return nil
}
