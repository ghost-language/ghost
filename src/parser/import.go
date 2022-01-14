package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) importStatement() ast.ExpressionNode {
	statement := &ast.Import{Token: parser.currentToken}

	if !parser.expectNextTokenIs(token.STRING) {
		return nil
	}

	statement.Path = &ast.String{Token: parser.currentToken, Value: parser.currentToken.Literal.(string)}

	parser.readToken()

	return statement
}
