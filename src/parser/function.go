package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) functionStatement() ast.ExpressionNode {
	expression := &ast.Function{Token: parser.currentToken}

	if !parser.expectNextType(token.LEFTPAREN) {
		expression.Name = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme}
		parser.readToken()
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

	if parser.nextTokenTypeIs(token.RIGHTPAREN) {
		parser.readToken()

		return defaults, parameters
	}

	// to do

	return defaults, parameters
}
