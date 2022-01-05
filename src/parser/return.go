package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) returnStatement() ast.StatementNode {
	statement := &ast.Return{Token: parser.currentToken}

	parser.readToken()

	statement.Value = parser.parseExpression(LOWEST)

	if parser.nextTokenIs(token.SEMICOLON) {
		parser.readToken()
	}

	return statement
}
