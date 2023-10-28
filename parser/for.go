package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) forExpression() ast.ExpressionNode {
	expression := &ast.For{Token: parser.currentToken}

	if !parser.expectNextTokenIs(token.LEFTPAREN) {
		return nil
	}

	parser.readToken()

	if !parser.currentTokenIs(token.IDENTIFIER) {
		return nil
	}

	if !parser.nextTokenIs(token.EQUAL) {
		return parser.forInExpression(expression)
	}

	expression.Identifier = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme}
	expression.Initializer = parser.assign()

	if expression.Initializer == nil {
		return nil
	}

	parser.readToken()

	expression.Condition = parser.parseExpression(LOWEST)

	if expression.Condition == nil {
		return nil
	}

	parser.readToken()
	parser.readToken()

	expression.Increment = parser.forIncrement()

	if expression.Increment == nil {
		return nil
	}

	if !parser.expectNextTokenIs(token.RIGHTPAREN) {
		return nil
	}

	if !parser.expectNextTokenIs(token.LEFTBRACE) {
		return nil
	}

	expression.Block = parser.blockStatement()

	return expression
}

func (parser *Parser) forInExpression(parent *ast.For) ast.ExpressionNode {
	expression := &ast.ForIn{Token: parent.Token}

	if !parser.currentTokenIs(token.IDENTIFIER) {
		return nil
	}

	value := ast.Identifier{Value: parser.currentToken.Lexeme}
	key := ast.Identifier{}

	parser.readToken()

	if parser.currentTokenIs(token.COMMA) {
		parser.readToken()

		if !parser.currentTokenIs(token.IDENTIFIER) {
			return nil
		}

		key = value
		value.Value = parser.currentToken.Lexeme

		parser.readToken()
	}

	expression.Key = &key
	expression.Value = &value

	if !parser.currentTokenIs(token.IN) {
		return nil
	}

	parser.readToken()

	expression.Iterable = parser.parseExpression(LOWEST)

	if !parser.expectNextTokenIs(token.RIGHTPAREN) {
		return nil
	}

	if !parser.expectNextTokenIs(token.LEFTBRACE) {
		return nil
	}

	expression.Block = parser.blockStatement()

	return expression
}

// forIncrement parses the increment expression of a for loop.
// It can be an assignment (x = x + 1), a postfix expression (x++), or an infix expression (x += 1).
func (parser *Parser) forIncrement() ast.ExpressionNode {
	if parser.currentTokenIs(token.RIGHTPAREN) {
		return nil
	}

	if parser.currentTokenIs(token.SEMICOLON) {
		parser.readToken()
		return nil
	}

	if parser.currentTokenIs(token.IDENTIFIER) && parser.nextTokenIs(token.EQUAL) {
		return parser.assign()
	}

	if parser.currentTokenIs(token.IDENTIFIER) && (parser.nextTokenIs(token.PLUSEQUAL) ||
		parser.nextTokenIs(token.MINUSEQUAL) ||
		parser.nextTokenIs(token.SLASHEQUAL) ||
		parser.nextTokenIs(token.STAREQUAL)) {
		identifier := parser.identifierLiteral()

		parser.readToken()

		return parser.compoundExpression(identifier)
	}

	if parser.currentTokenIs(token.IDENTIFIER) && (parser.nextTokenIs(token.PLUSPLUS) || parser.nextTokenIs(token.MINUSMINUS)) {
		parser.readToken()

		return parser.postfixExpression()
	}

	return nil
}
