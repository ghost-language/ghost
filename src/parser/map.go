package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) mapLiteral() ast.ExpressionNode {
	mapLiteral := &ast.Map{Token: parser.currentToken}
	mapLiteral.Pairs = make(map[ast.ExpressionNode]ast.ExpressionNode)

	for !parser.nextTokenTypeIs(token.RIGHTBRACE) {
		parser.readToken()

		key := parser.parseExpression(LOWEST)

		if !parser.expectNextType(token.COLON) {
			return nil
		}

		parser.readToken()

		value := parser.parseExpression(LOWEST)

		mapLiteral.Pairs[key] = value

		if !parser.nextTokenTypeIs(token.RIGHTBRACE) && !parser.expectNextType(token.COMMA) {
			return nil
		}
	}

	if !parser.expectNextType(token.RIGHTBRACE) {
		return nil
	}

	return mapLiteral
}
