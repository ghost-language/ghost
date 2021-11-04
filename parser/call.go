package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) callExpression(callee ast.ExpressionNode) ast.ExpressionNode {
	call := &ast.Call{Token: parser.current(), Callee: callee}

	call.Arguments = parser.callArguments()

	return call
}

func (parser *Parser) callArguments() []ast.ExpressionNode {
	args := []ast.ExpressionNode{}

	if parser.checkNext(token.RIGHTPAREN) {
		parser.advance()
		return args
	}

	parser.advance()
	args = append(args, parser.parseExpression(LOWEST))

	for parser.checkNext(token.COMMA) {
		parser.advance()
		parser.advance()
		args = append(args, parser.parseExpression(LOWEST))
	}

	if !parser.checkNext(token.RIGHTPAREN) {
		return nil
	}

	return args
}
