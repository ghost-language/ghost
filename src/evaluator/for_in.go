package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"github.com/shopspring/decimal"
)

func evaluateForIn(node *ast.ForIn, scope *object.Scope) object.Object {
	iterable := Evaluate(node.Iterable, scope)

	if isError(iterable) {
		return iterable
	}

	existingKey, keyExisted := scope.Environment.Get(node.Key.Value)
	existingValue, valueExisted := scope.Environment.Get(node.Value.Value)

	defer func() {
		if keyExisted {
			scope.Environment.Set(node.Key.Value, existingKey)
		} else {
			scope.Environment.Delete(node.Key.Value)
		}

		if valueExisted {
			scope.Environment.Set(node.Value.Value, existingValue)
		} else {
			scope.Environment.Delete(node.Value.Value)
		}
	}()

	switch obj := iterable.(type) {
	case *object.List:
		for k, v := range obj.Elements {
			scope.Environment.Set(node.Key.Value, &object.Number{Value: decimal.NewFromInt(int64(k))})
			scope.Environment.Set(node.Value.Value, v)

			block := Evaluate(node.Block, scope)

			if isError(block) {
				return block
			}
		}

		return nil
	case *object.Map:
		for _, pair := range obj.Pairs {
			scope.Environment.Set(node.Key.Value, pair.Key)
			scope.Environment.Set(node.Value.Value, pair.Value)

			block := Evaluate(node.Block, scope)

			if isError(block) {
				return block
			}
		}

		return nil
	}

	return newError("%d:%d:%s: runtime error: unusable as for loop: %T", node.Token.Line, node.Token.Column, node.Token.File, iterable)
}
