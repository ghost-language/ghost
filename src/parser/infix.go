package parser

import (
	"ghostlang.org/x/ghost/ast"
)

func (parser *Parser) infixExpression(left ast.ExpressionNode) ast.ExpressionNode {
	infix := &ast.Infix{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Lexeme,
		Left:     left,
	}

	precedence := parser.currentTokenPrecedence()

	parser.readToken()

	infix.Right = parser.parseExpression(precedence)

	return infix
}
