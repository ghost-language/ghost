package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"github.com/shopspring/decimal"
)

func evaluateForIn(node *ast.ForIn, env *object.Environment) object.Object {
	iterable := Evaluate(node.Iterable, env)

	existingKey, keyExisted := env.Get(node.Key.Value)
	existingValue, valueExisted := env.Get(node.Value.Value)

	defer func() {
		if keyExisted {
			env.Set(node.Key.Value, existingKey)
		} else {
			env.Delete(node.Key.Value)
		}

		if valueExisted {
			env.Set(node.Value.Value, existingValue)
		} else {
			env.Delete(node.Value.Value)
		}
	}()

	switch obj := iterable.(type) {
	case *object.List:
		for k, v := range obj.Elements {
			env.Set(node.Key.Value, &object.Number{Value: decimal.NewFromInt(int64(k))})
			env.Set(node.Value.Value, v)

			block := Evaluate(node.Block, env)

			if isError(block) {
				return block
			}
		}

		return nil
	case *object.Map:
		for _, pair := range obj.Pairs {
			env.Set(node.Key.Value, pair.Key)
			env.Set(node.Value.Value, pair.Value)

			block := Evaluate(node.Block, env)

			if isError(block) {
				return block
			}
		}

		return nil
	}

	return newError("%d:%d: runtime error: unusable as for loop: %T", node.Token.Line, node.Token.Column, iterable)
}
