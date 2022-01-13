package parser

import (
	"ghostlang.org/x/ghost/ast"
)

func (parser *Parser) compoundExpression(left ast.ExpressionNode) ast.ExpressionNode {
	compound := &ast.Compound{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Lexeme,
		Left:     left,
	}

	precedence := parser.currentTokenPrecedence()

	parser.readToken()

	compound.Right = parser.parseExpression(precedence)

	return compound
}
