package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) mapLiteral() ast.ExpressionNode {
	mapLiteral := &ast.Map{Token: parser.currentToken}
	mapLiteral.Pairs = []ast.MapEntry{}

	for !parser.nextTokenIs(token.RIGHTBRACE) {
		parser.readToken()

		key := parser.parseExpression(LOWEST)

		identifier, isIdentifier := key.(*ast.Identifier)

		if isIdentifier && (parser.nextTokenIs(token.COMMA) || parser.nextTokenIs(token.RIGHTBRACE)) {
			mapLiteral.Pairs = append(mapLiteral.Pairs, ast.MapEntry{
				Key:   key,
				Value: &ast.Identifier{Token: identifier.Token, Value: identifier.Value},
			})

			if !parser.nextTokenIs(token.RIGHTBRACE) && !parser.expectNextTokenIs(token.COMMA) {
				return nil
			}

			continue
		}

		if !parser.expectNextTokenIs(token.COLON) {
			return nil
		}

		parser.readToken()

		value := parser.parseExpression(LOWEST)

		mapLiteral.Pairs = append(mapLiteral.Pairs, ast.MapEntry{Key: key, Value: value})

		if !parser.nextTokenIs(token.RIGHTBRACE) && !parser.expectNextTokenIs(token.COMMA) {
			return nil
		}
	}

	if !parser.expectNextTokenIs(token.RIGHTBRACE) {
		return nil
	}

	return mapLiteral
}
