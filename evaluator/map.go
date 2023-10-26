package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
)

func evaluateMap(node *ast.Map, scope *object.Scope) object.Object {
	pairs := make(map[object.MapKey]object.MapPair)

	for keyNode, valueNode := range node.Pairs {
		// if keyNode is an identifier, convert it to a string
		identifier, ok := keyNode.(*ast.Identifier)

		if ok {
			keyNode = &ast.String{
				Token: identifier.Token,
				Value: identifier.Value,
			}
		}

		key := Evaluate(keyNode, scope)

		if isError(key) {
			return key
		}

		mapKey, ok := key.(object.Mappable)

		if !ok {
			return newError("%d:%d:%s: runtime error: unusable as map key: %s", node.Token.Line, node.Token.Column, node.Token.File, key.Type())
		}

		value := Evaluate(valueNode, scope)

		if isError(value) {
			return value
		}

		hashed := mapKey.MapKey()

		pairs[hashed] = object.MapPair{Key: key, Value: value}
	}

	return &object.Map{Pairs: pairs}
}
