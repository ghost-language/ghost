package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
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
			return object.NewError(fault.Type, node.Token, "%s cannot be used as a map key", object.TypeName(key)).
				WithHelp("a map key has to be a string, a number, or a boolean")
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
