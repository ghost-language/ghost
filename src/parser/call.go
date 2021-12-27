package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) callExpression(callee ast.ExpressionNode) ast.ExpressionNode {
	call := &ast.Call{Token: parser.currentToken, Callee: callee}

	call.Arguments = parser.callArguments()

	return call
}

func (parser *Parser) callArguments() []ast.ExpressionNode {
	args := []ast.ExpressionNode{}

	if parser.nextTokenTypeIs(token.RIGHTPAREN) {
		return args
	}

	parser.readToken()

	args = append(args, parser.parseExpression(LOWEST))

	for parser.nextTokenTypeIs(token.COMMA) {
		parser.readToken()
		args = append(args, parser.parseExpression(LOWEST))
	}

	if !parser.expectNextType(token.RIGHTPAREN) {
		return nil
	}

	return args
}
