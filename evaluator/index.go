package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
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

	if left == nil || index == nil {
		return object.NewError(fault.Type, node.Token, "cannot index a null value")
	}

	switch {
	case left.Type() == object.STRING && index.Type() == object.NUMBER:
		return evaluateStringIndex(node, left, index)
	case left.Type() == object.LIST && index.Type() == object.NUMBER:
		return evaluateListIndex(node, left, index)
	case left.Type() == object.MAP:
		return evaluateMapIndex(node, left, index)
	default:
		return indexTypeError(node, left, index)
	}
}

func evaluateListIndex(node *ast.Index, left, index object.Object) object.Object {
	list := left.(*object.List)
	idx := index.(*object.Number).Int64()
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
		return object.NewError(fault.Type, node.Token, "%s cannot be used as a map key", object.TypeName(index)).
			WithHelp("a map key has to be a string, a number, or a boolean")
	}

	pair, ok := mapObject.Pairs[key.MapKey()]

	if !ok {
		return value.NULL
	}

	return pair.Value
}

func evaluateStringIndex(node *ast.Index, left, index object.Object) object.Object {
	str := left.(*object.String)
	idx := index.(*object.Number).Int64()
	max := int64(len(str.Value) - 1)

	if idx < 0 || idx > max {
		return value.NULL
	}

	return &object.String{Value: string([]rune(str.Value)[idx])}
}
