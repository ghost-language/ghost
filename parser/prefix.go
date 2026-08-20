package parser

import "ghostlang.org/x/ghost/ast"

func (parser *Parser) prefixExpression() ast.ExpressionNode {
	prefix := &ast.Prefix{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Type,
	}

	parser.readToken()

	prefix.Right = parser.parseExpression(PREFIX)

	return prefix
}
