package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

func evaluatePostfix(node *ast.Postfix, scope *object.Scope) object.Object {
	name := node.Token.Lexeme

	number, err := readCounter(node, scope, name)

	if err != nil {
		return err
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

	scope.Environment.Set(name, updated)

	return updated
}

// readCounter reads the variable a `++` or `--` is applied to. Both operators
// fail the same two ways, so they ask the same question and report the same
// answers rather than each spelling them out.
func readCounter(node *ast.Postfix, scope *object.Scope, name string) (*object.Number, *object.Error) {
	current, ok := scope.Environment.Get(name)

	if !ok {
		return nil, undefined(node.Token, name, scope)
	}

	number, ok := current.(*object.Number)

	if !ok {
		return nil, object.NewError(fault.Type, node.Token, "cannot use `%s` on `%s`, which is a %s, not a number", node.Operator, name, object.TypeName(current))
	}

	return number, nil
}
