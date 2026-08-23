package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

func evaluatePostfix(node *ast.Postfix, scope *object.Scope) object.Object {
	current := Evaluate(node.Left, scope)

	if isError(current) {
		return current
	}

	number, ok := current.(*object.Number)

	if !ok {
		return object.NewError(fault.Type, node.Token, "cannot use `%s` on %s", node.Operator, object.TypeName(current))
	}

	var updated object.Object

	switch node.Operator {
	case token.PLUSPLUS:
		updated = number.Increment()
	case token.MINUSMINUS:
		updated = number.Decrement()
	default:
		return object.NewError(fault.Internal, node.Token, "`%s` was parsed as a postfix operator but has no behaviour behind it", node.Operator)
	}

	switch target := node.Left.(type) {
	case *ast.Identifier:
		scope.Environment.Set(target.Value, updated)
	case *ast.Index:
		if result := evaluateIndexAssignment(target, updated, scope); isError(result) {
			return result
		}
	case *ast.Property:
		if result := evaluatePropertyAssignment(target, updated, scope); isError(result) {
			return result
		}
	default:
		return object.NewError(fault.Syntax, node.Token, "cannot assign to this expression with `%s`", node.Operator).
			WithHelp("only a variable, a list or map entry, or a property can be assigned to")
	}

	return updated
}
