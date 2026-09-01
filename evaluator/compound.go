package evaluator

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

// compoundOperators maps each compound assignment operator to the binary
// operator it applies, e.g. `+=` performs `+`.
var compoundOperators = map[token.Type]token.Type{
	token.PLUSEQUAL:  token.PLUS,
	token.MINUSEQUAL: token.MINUS,
	token.STAREQUAL:  token.STAR,
	token.SLASHEQUAL: token.SLASH,
}

func evaluateCompound(node *ast.Compound, scope *object.Scope) object.Object {
	operator, ok := compoundOperators[node.Operator]

	if !ok {
		return object.NewError(fault.Internal, node.Token, "`%s` was parsed as a compound assignment but has no operator behind it", node.Operator)
	}

	infix := &ast.Infix{
		Token:    node.Token,
		Left:     node.Left,
		Operator: operator,
		Right:    node.Right,
	}

	value := evaluateInfix(infix, scope)

	if isError(value) {
		return value
	}

	switch target := node.Left.(type) {
	case *ast.Identifier:
		bind(scope, target.Value, value)
	case *ast.Index:
		return evaluateIndexAssignment(target, value, scope)
	case *ast.Property:
		return evaluatePropertyAssignment(target, value, scope)
	default:
		return object.NewError(fault.Syntax, node.Token, "cannot assign to this expression with `%s`", node.Operator).
			WithHelp("only a variable, a list or map entry, or a property can be assigned to")
	}

	return nil
}
