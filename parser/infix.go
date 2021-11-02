package parser

import (
	"ghostlang.org/x/ghost/ast"
)

func (parser *Parser) infixExpression(left ast.ExpressionNode) ast.ExpressionNode {
	infix := &ast.Infix{
		Token:    parser.current(),
		Operator: parser.current().Lexeme,
		Left:     left,
	}

	parser.advance()

	precedence := parser.currentPrecedence()

	infix.Right = parser.parseExpression(precedence)

	return infix
}
