package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) callExpression(callee ast.ExpressionNode) ast.ExpressionNode {
	call := &ast.Call{Token: parser.currentToken, Callee: callee}

	call.Arguments = parser.parseExpressionList(token.RIGHTPAREN)

	return call
}
