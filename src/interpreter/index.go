package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateIndex(node *ast.Index, env *object.Environment) object.Object {
	left := Evaluate(node.Left, env)

	if isError(left) {
		return left
	}

	index := Evaluate(node.Index, env)

	if isError(index) {
		return index
	}

	switch {
	case left.Type() == object.STRING && index.Type() == object.NUMBER:
		return evaluateStringIndex(left, index)
	case left.Type() == object.LIST && index.Type() == object.NUMBER:
		return evaluateListIndex(left, index)
	case left.Type() == object.MAP:
		return evaluateMapIndex(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evaluateListIndex(left, index object.Object) object.Object {
	list := left.(*object.List)
	idx := index.(*object.Number).Value.IntPart()
	max := int64(len(list.Elements) - 1)

	if idx < 0 || idx > max {
		return value.NULL
	}

	return list.Elements[idx]
}

func evaluateMapIndex(left, index object.Object) object.Object {
	mapObject := left.(*object.Map)

	key, ok := index.(object.Mappable)

	if !ok {
		err := newError("unusuable as map key: %s", index.Type())

		return err
	}

	pair, ok := mapObject.Pairs[key.MapKey()]

	if !ok {
		return value.NULL
	}

	return pair.Value
}

func evaluateStringIndex(left, index object.Object) object.Object {
	str := left.(*object.String)
	idx := index.(*object.Number).Value.IntPart()
	max := int64(len(str.Value) - 1)

	if idx < 0 || idx > max {
		return value.NULL
	}

	return &object.String{Value: string([]rune(str.Value)[idx])}
}
