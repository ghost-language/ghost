package evaluator

import (
	"ghostlang.org/x/ghost/ast"
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
		return newError("%d:%d:%s: runtime error: unknown operator: %s", node.Token.Line, node.Token.Column, node.Token.File, node.Operator)
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
		scope.Environment.Set(target.Value, value)
	case *ast.Index:
		return evaluateIndexAssignment(target, value, scope)
	case *ast.Property:
		return evaluatePropertyAssignment(target, value, scope)
	default:
		return newError("%d:%d:%s: runtime error: invalid compound assignment target: %T", node.Token.Line, node.Token.Column, node.Token.File, node.Left)
	}

	return nil
}
