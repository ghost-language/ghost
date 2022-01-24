package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) switchStatement() ast.ExpressionNode {
	expression := &ast.Switch{Token: parser.currentToken}

	if !parser.expectNextTokenIs(token.LEFTPAREN) {
		return nil
	}

	parser.readToken()

	expression.Value = parser.parseExpression(LOWEST)

	if expression.Value == nil {
		return nil
	}

	if !parser.expectNextTokenIs(token.RIGHTPAREN) {
		return nil
	}

	// block of cases
	if !parser.expectNextTokenIs(token.LEFTBRACE) {
		return nil
	}

	parser.readToken()

	for !parser.currentTokenIs(token.RIGHTBRACE) {
		// check for EOF
		//

		switchCase := &ast.Case{Token: parser.currentToken}

		if parser.currentTokenIs(token.CASE) {
			// read "case"
			parser.readToken()

			// A switch case can contain multiple "values"
			switchCase.Value = append(switchCase.Value, parser.parseExpression(LOWEST))

			for parser.nextTokenIs(token.COMMA) {
				// read the comma
				parser.readToken()

				// setup the expression
				parser.readToken()

				switchCase.Value = append(switchCase.Value, parser.parseExpression(LOWEST))
			}
		}

		if !parser.expectNextTokenIs(token.COLON) {
			return nil
		}

		if !parser.expectNextTokenIs(token.LEFTBRACE) {
			return nil
		}

		// parse the block
		switchCase.Body = parser.blockStatement()

		if !parser.currentTokenIs(token.RIGHTBRACE) {
			return nil
		}

		parser.readToken()

		expression.Cases = append(expression.Cases, switchCase)
	}

	if !parser.currentTokenIs(token.RIGHTBRACE) {
		return nil
	}

	return expression
}
