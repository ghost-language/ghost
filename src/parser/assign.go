package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) assign() ast.StatementNode {
	if parser.currentTokenIs(token.IDENTIFIER) && parser.nextTokenIs(token.ASSIGN) {
		statement := &ast.Assign{
			Name:  &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme},
			Token: parser.nextToken,
		}

		if !parser.expectNextTokenIs(token.ASSIGN) {
			return nil
		}

		parser.readToken()

		statement.Value = parser.parseExpression(LOWEST)

		if parser.nextTokenIs(token.SEMICOLON) {
			parser.readToken()
		}

		return statement
	}

	return nil
}
