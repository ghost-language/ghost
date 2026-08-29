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

// forIncrement parses the increment clause of a for loop. A bare identifier
// followed by `=` goes through assign() the same as any other assignment
// statement, since plain `=` has no infix parser of its own (§ assign.go).
// Everything else — a postfix `x++`/`x--`, a compound `x += 1`, or either of
// those targeting a property or index (`obj.count++`, `list[0] += 1`) —
// already has an infix parser registered in the normal expression table, so
// parseExpression handles it the same way it would anywhere else in a
// program; the increment clause has no business special-casing identifiers
// only.
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

	return parser.parseExpression(LOWEST)
}
