package parser

import (
	"ghostlang.org/x/ghost/ast"
)

func (parser *Parser) expression(precedence int) ast.ExpressionNode {
	postfix := parser.postfixParserFns[parser.current().Type]

	if postfix != nil {
		return postfix()
	}

	prefix := parser.prefixParserFns[parser.current().Type]

	if prefix == nil {
		return nil
	}

	leftExpression := prefix()

	for !parser.isAtEnd() && precedence < parser.nextPrecedence() {
		infix := parser.infixParserFns[parser.next().Type]

		if infix == nil {
			return leftExpression
		}

		leftExpression = infix(leftExpression)
	}

	return leftExpression
}
