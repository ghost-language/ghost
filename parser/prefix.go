package parser

import "ghostlang.org/x/ghost/ast"

func (parser *Parser) prefixExpression() ast.ExpressionNode {
	prefix := &ast.Prefix{
		Token:    parser.current(),
		Operator: parser.current().Lexeme,
	}

	parser.advance()

	prefix.Right = parser.parseExpression(PREFIX)

	return prefix
}
