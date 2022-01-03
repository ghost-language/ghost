package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateIndex(node *ast.Index, env *object.Environment) (object.Object, bool) {
	left, ok := Evaluate(node.Left, env)

	if !ok {
		return nil, false
	}

	index, ok := Evaluate(node.Index, env)

	if !ok {
		return nil, false
	}

	switch {
	case left.Type() == object.LIST && index.Type() == object.NUMBER:
		return evaluateListIndex(left, index)
	case left.Type() == object.MAP:
		return evaluateMapIndex(left, index)
	default:
		err := newError("index operator not supported: %s", left.Type())

		return err, false
	}
}

func evaluateListIndex(left, index object.Object) (object.Object, bool) {
	list := left.(*object.List)
	idx := index.(*object.Number).Value.IntPart()
	max := int64(len(list.Elements) - 1)

	if idx < 0 || idx > max {
		return value.NULL, true
	}

	return list.Elements[idx], true
}

func evaluateMapIndex(left, index object.Object) (object.Object, bool) {
	mapObject := left.(*object.Map)

	key, ok := index.(object.Mappable)

	if !ok {
		err := newError("unusuable as map key: %s", index.Type())

		return err, false
	}

	pair, ok := mapObject.Pairs[key.MapKey()]

	if !ok {
		return value.NULL, true
	}

	return pair.Value, true
}
