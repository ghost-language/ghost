package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

// Arithmetic on lists is elementwise, and broadcasts: a number spreads across
// every element, and a shorter shape stretches across a longer one. The rules
// are object.Broadcast's, which is also what the math module's methods use, so
// `a + b` and `math.add(a, b)` are the same operation reached two ways.
//
//	[1, 2, 3] * 2                 // [2, 4, 6]
//	[1, 2, 3] + [10, 20, 30]      // [11, 22, 33]
//	[[1, 2], [3, 4]] + [10, 20]   // [[11, 22], [13, 24]]
//
// Joining two lists end to end is a different operation with a different
// meaning, and has a method rather than an operator: `first.concat(second)`.
func evaluateListInfix(node *ast.Infix, left object.Object, right object.Object) object.Object {
	operation, ok := listOperations[node.Operator]

	if !ok {
		return newError("%d:%d:%s: runtime error: unknown operator: %s %s %s", node.Token.Line, node.Token.Column, node.Token.File, left.Type(), node.Operator, right.Type())
	}

	result, fault := object.Broadcast([]object.Object{left, right}, func(values []*object.Number) object.Object {
		return operation(node, values[0], values[1])
	})

	if fault != nil {
		return newError("%d:%d:%s: runtime error: cannot evaluate %s %s %s: %s", node.Token.Line, node.Token.Column, node.Token.File, left.Type(), node.Operator, right.Type(), fault.Reason)
	}

	return result
}

// listOperations is the set of operators that mean something elementwise. The
// comparisons are deliberately absent: `==` is handled separately, as a
// structural comparison of whole lists, and an ordering between two lists has
// no reading obvious enough to pick one.
var listOperations = map[token.Type]func(node *ast.Infix, left *object.Number, right *object.Number) object.Object{
	token.PLUS:  func(node *ast.Infix, left *object.Number, right *object.Number) object.Object { return left.Add(right) },
	token.MINUS: func(node *ast.Infix, left *object.Number, right *object.Number) object.Object { return left.Sub(right) },
	token.STAR:  func(node *ast.Infix, left *object.Number, right *object.Number) object.Object { return left.Mul(right) },
	token.SLASH: func(node *ast.Infix, left *object.Number, right *object.Number) object.Object {
		if right.IsZero() {
			return newError("%d:%d:%s: runtime error: division by zero", node.Token.Line, node.Token.Column, node.Token.File)
		}

		return left.Div(right)
	},
	token.PERCENT: func(node *ast.Infix, left *object.Number, right *object.Number) object.Object {
		if right.IsZero() {
			return newError("%d:%d:%s: runtime error: division by zero", node.Token.Line, node.Token.Column, node.Token.File)
		}

		return left.Mod(right)
	},
}

// listsEqual compares two lists by their contents rather than by identity, to
// any depth. It is what `==` between lists means, and mirrors the comparisons
// the type-specific infix evaluators make for the values inside.
func listsEqual(left *object.List, right *object.List) bool {
	if len(left.Elements) != len(right.Elements) {
		return false
	}

	for index, element := range left.Elements {
		if !valuesEqual(element, right.Elements[index]) {
			return false
		}
	}

	return true
}

func valuesEqual(left object.Object, right object.Object) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}

	if left.Type() != right.Type() {
		return false
	}

	switch left := left.(type) {
	case *object.Number:
		return left.Equal(right.(*object.Number))
	case *object.String:
		return left.Value == right.(*object.String).Value
	case *object.Boolean:
		return left.Value == right.(*object.Boolean).Value
	case *object.Null:
		return true
	case *object.List:
		return listsEqual(left, right.(*object.List))
	}

	// Everything else - instances, functions, maps - compares by identity, the
	// same rule `==` applies to instances outside a list.
	return left == right
}
