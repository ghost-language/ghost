package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) assign() ast.StatementNode {
	statement := &ast.Assign{Token: parser.currentToken}

	if parser.currentTokenIs(token.IDENTIFIER) {
		statement.Name = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme}
	} else if parser.currentTokenIs(token.ASSIGN) {
		// index or property assignment

		if parser.previousIndex != nil {
			// foo["bar"] := true
			statement.Name = parser.previousIndex

			parser.readToken()

			statement.Value = parser.parseExpression(LOWEST)

			parser.previousIndex = nil

			if parser.nextTokenIs(token.SEMICOLON) {
				parser.readToken()
			}

			return statement
		}

		if parser.previousProperty != nil {
			// foo.bar := true
			statement.Name = parser.previousProperty

			parser.readToken()

			statement.Value = parser.parseExpression(LOWEST)

			parser.previousProperty = nil

			if parser.nextTokenIs(token.SEMICOLON) {
				parser.readToken()
			}

			return statement
		}
	}

	if !parser.nextTokenIs(token.ASSIGN) {
		return nil
	}

	parser.readToken()
	statement.Token = parser.currentToken
	parser.readToken()

	statement.Value = parser.parseExpression(LOWEST)

	if parser.nextTokenIs(token.SEMICOLON) {
		parser.readToken()
	}

	return statement
}
