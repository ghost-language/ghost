package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) mapLiteral() ast.ExpressionNode {
	mapLiteral := &ast.Map{Token: parser.currentToken}
	mapLiteral.Pairs = make(map[ast.ExpressionNode]ast.ExpressionNode)

	for !parser.nextTokenIs(token.RIGHTBRACE) {
		parser.readToken()

		key := parser.parseExpression(LOWEST)

		if !parser.expectNextTokenIs(token.COLON) {
			return nil
		}

		parser.readToken()

		value := parser.parseExpression(LOWEST)

		mapLiteral.Pairs[key] = value

		log.Debug("next token: %s", parser.nextToken.Lexeme)

		if !parser.currentTokenIs(token.COMMA) || !parser.nextTokenIs(token.RIGHTBRACE) && !parser.expectNextTokenIs(token.COMMA) {
			return nil
		}
	}

	if !parser.expectNextTokenIs(token.RIGHTBRACE) {
		return nil
	}

	return mapLiteral
}
