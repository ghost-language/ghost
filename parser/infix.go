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

	precedence := parser.peekPrecedence()

	parser.advance()

	infix.Right = parser.parseExpression(precedence)

	return infix
}
