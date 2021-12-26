package parser

import (
	"ghostlang.org/x/ghost/ast"
)

func (parser *Parser) parseExpression(precedence int) ast.ExpressionNode {
	prefix := parser.prefixParserFns[parser.currentToken.Type]

	if prefix == nil {
		return nil
	}

	leftExpression := prefix()

	for precedence < parser.nextTokenPrecedence() {
		infix := parser.infixParserFns[parser.nextToken.Type]

		if infix == nil {
			return leftExpression
		}

		parser.readToken()

		leftExpression = infix(leftExpression)
	}

	return leftExpression
}
