package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateMap(node *ast.Map, env *object.Environment) (object.Object, bool) {
	pairs := make(map[object.MapKey]object.MapPair)

	for keyNode, valueNode := range node.Pairs {
		key, ok := Evaluate(keyNode, env)

		if !ok {
			return nil, false
		}

		mapKey, ok := key.(object.Mappable)

		if !ok {
			err := newError("unusable as map key: %s", key.Type())

			return err, false
		}

		value, ok := Evaluate(valueNode, env)

		if !ok {
			return nil, false
		}

		hashed := mapKey.MapKey()

		pairs[hashed] = object.MapPair{Key: key, Value: value}
	}

	return &object.Map{Pairs: pairs}, true
}
