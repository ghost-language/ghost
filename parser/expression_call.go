package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (p *Parser) parseCallExpression(callable ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Token: p.currentToken, Callable: callable}
	expression.Arguments = p.parseExpressionList(token.RPAREN)

	return expression
}
