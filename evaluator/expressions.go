package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
)

// evaluateExpressions evaluates a call's arguments or a list literal's
// elements - the two places `...expr` (ast.Spread, §12) is meaningful. A
// spread expands the list it evaluates to in place, rather than nesting that
// list as a single argument or element, so `f(...a, b)` passes each of a's
// elements followed by b, and `[...a, b]` builds one flat list the same way.
func evaluateExpressions(expressions []ast.ExpressionNode, scope *object.Scope) []object.Object {
	var result []object.Object

	for _, expression := range expressions {
		if spread, ok := expression.(*ast.Spread); ok {
			evaluated := Evaluate(spread.Value, scope)

			if isError(evaluated) {
				return []object.Object{evaluated}
			}

			list, ok := evaluated.(*object.List)

			if !ok {
				err := object.NewError(fault.Type, spread.Token, "cannot spread %s, only a list", object.TypeName(evaluated)).
					WithHelp("`...` expands a list's elements in place, as in `f(...list)` or `[...list, 1]`")

				return []object.Object{err}
			}

			result = append(result, list.Elements...)

			continue
		}

		evaluated := Evaluate(expression, scope)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}
