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

		// Read the IDENTIFIER and ASSIGN tokens
		parser.readToken()
		parser.readToken()

		statement.Value = parser.parseExpression(LOWEST)

		return statement
	}

	return nil
}
