package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) functionStatement() ast.ExpressionNode {
	expression := &ast.Function{Token: parser.currentToken}

	if !parser.nextTokenIs(token.LEFTPAREN) {
		parser.readToken()
		expression.Name = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme}
	}

	if !parser.expectNextType(token.LEFTPAREN) {
		return nil
	}

	expression.Defaults, expression.Parameters = parser.functionParameters()

	expression.Body = parser.blockStatement()

	return expression
}

func (parser *Parser) functionParameters() (map[string]ast.ExpressionNode, []*ast.Identifier) {
	defaults := make(map[string]ast.ExpressionNode)
	parameters := []*ast.Identifier{}

	if parser.nextTokenIs(token.RIGHTPAREN) {
		parser.readToken()

		return defaults, parameters
	}

	// to do

	return defaults, parameters
}
