package parser

import (
	"ghostlang.org/ghost/ast"
	"ghostlang.org/ghost/token"
)

func (p *Parser) parseMapLiteral() ast.Expression {
	mapLiteral := &ast.MapLiteral{Token: p.currentToken}
	mapLiteral.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()

		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()

		value := p.parseExpression(LOWEST)
		mapLiteral.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return mapLiteral
}
