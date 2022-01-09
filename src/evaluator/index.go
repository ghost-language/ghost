package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateIndex(node *ast.Index, scope *object.Scope) object.Object {
	left := Evaluate(node.Left, scope)

	if isError(left) {
		return left
	}

	index := Evaluate(node.Index, scope)

	if isError(index) {
		return index
	}

	switch {
	case left.Type() == object.STRING && index.Type() == object.NUMBER:
		return evaluateStringIndex(node, left, index)
	case left.Type() == object.LIST && index.Type() == object.NUMBER:
		return evaluateListIndex(node, left, index)
	case left.Type() == object.MAP:
		return evaluateMapIndex(node, left, index)
	default:
		return newError("%d:%d: runtime error: index operator not supported: %s", node.Token.Line, node.Token.Column, left.Type())
	}
}

func evaluateListIndex(node *ast.Index, left, index object.Object) object.Object {
	list := left.(*object.List)
	idx := index.(*object.Number).Value.IntPart()
	max := int64(len(list.Elements) - 1)

	if idx < 0 || idx > max {
		return value.NULL
	}

	return list.Elements[idx]
}

func evaluateMapIndex(node *ast.Index, left, index object.Object) object.Object {
	mapObject := left.(*object.Map)

	key, ok := index.(object.Mappable)

	if !ok {
		return newError("%d:%d: runtime error: unusable as map key: %s", node.Token.Line, node.Token.Column, index.Type())
	}

	pair, ok := mapObject.Pairs[key.MapKey()]

	if !ok {
		return value.NULL
	}

	return pair.Value
}

func evaluateStringIndex(node *ast.Index, left, index object.Object) object.Object {
	str := left.(*object.String)
	idx := index.(*object.Number).Value.IntPart()
	max := int64(len(str.Value) - 1)

	if idx < 0 || idx > max {
		return value.NULL
	}

	return &object.String{Value: string([]rune(str.Value)[idx])}
}
