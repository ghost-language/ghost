package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

// evaluateDateInfix covers comparing two dates - the one family of operations
// for which two dates have an obvious, single reading. See object.Date's doc
// comment for why arithmetic between two dates is deliberately not here.
func evaluateDateInfix(node *ast.Infix, left object.Object, right object.Object) object.Object {
	leftDate := left.(*object.Date)
	rightDate := right.(*object.Date)

	switch node.Operator {
	case token.LESS:
		return toBooleanValue(leftDate.Time.Before(rightDate.Time))
	case token.LESSEQUAL:
		return toBooleanValue(!leftDate.Time.After(rightDate.Time))
	case token.GREATER:
		return toBooleanValue(leftDate.Time.After(rightDate.Time))
	case token.GREATEREQUAL:
		return toBooleanValue(!leftDate.Time.Before(rightDate.Time))
	case token.EQUALEQUAL:
		return toBooleanValue(leftDate.Time.Equal(rightDate.Time))
	case token.BANGEQUAL:
		return toBooleanValue(!leftDate.Time.Equal(rightDate.Time))
	}

	return object.NewError(fault.Type, node.Token, "cannot use `%s` between two dates", node.Operator).
		WithHelp("date arithmetic goes through the date module - addDays, differenceInDays, and the rest")
}
