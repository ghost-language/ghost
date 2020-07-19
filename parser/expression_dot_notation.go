package parser

import (
	"ghostlang.org/ghost/ast"
	"ghostlang.org/ghost/token"
)

func (p *Parser) parseDotNotationExpression(expression ast.Expression) ast.Expression {
	p.expectPeek(token.IDENTIFIER)

	index := &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}

	return &ast.IndexExpression{Left: expression, Index: index}
}
