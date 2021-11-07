package parser

import (
	"ghostlang.org/x/ghost/ast"
)

func (parser *Parser) parseExpression(precedence int) ast.ExpressionNode {
	postfix := parser.postfixParserFns[parser.peek().Type]

	if postfix != nil {
		return postfix()
	}

	prefix := parser.prefixParserFns[parser.peek().Type]

	if prefix == nil {
		return nil
	}

	leftExpression := prefix()

	for precedence < parser.nextPrecedence() {
		infix := parser.infixParserFns[parser.next().Type]

		if infix == nil {
			return leftExpression
		}

		parser.advance()

		leftExpression = infix(leftExpression)
	}

	return leftExpression
}
