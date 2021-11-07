package parser

import (
	"ghostlang.org/x/ghost/ast"
)

func (parser *Parser) infixExpression(left ast.ExpressionNode) ast.ExpressionNode {
	infix := &ast.Infix{
		Token:    parser.peek(),
		Operator: parser.peek().Lexeme,
		Left:     left,
	}

	parser.advance()

	precedence := parser.peekPrecedence()

	infix.Right = parser.parseExpression(precedence)

	return infix
}
