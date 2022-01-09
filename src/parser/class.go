package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) classStatement() ast.ExpressionNode {
	class := &ast.Class{Token: parser.currentToken}

	parser.readToken()

	class.Name = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme}

	if !parser.expectNextTokenIs(token.LEFTBRACE) {
		return nil
	}

	class.Body = parser.blockStatement()

	return class
}
