package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) returnStatement() ast.StatementNode {
	statement := &ast.Return{Token: parser.currentToken}

	if parser.nextTokenIs(token.SEMICOLON) {
		statement.Value = &ast.Null{Token: parser.currentToken}
	} else if parser.nextTokenIs(token.RIGHTBRACE) || parser.nextTokenIs(token.EOF) {
		statement.Value = &ast.Null{Token: parser.currentToken}
	} else {
		parser.readToken()

		statement.Value = parser.parseExpression(LOWEST)
	}

	if parser.nextTokenIs(token.SEMICOLON) {
		parser.readToken()
	}

	return statement
}
