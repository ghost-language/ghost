package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) functionStatement() ast.ExpressionNode {
	expression := &ast.Function{Token: parser.advance()}

	if !parser.check(token.LEFTPAREN) {
		expression.Name = &ast.Identifier{Token: parser.peek(), Value: parser.peek().Lexeme}
		parser.advance()
	}

	parser.consume(token.LEFTPAREN, "Expect '(' after 'function'. got=%s", parser.peek().Type)

	expression.Defaults, expression.Parameters = parser.functionParameters()

	expression.Body = parser.blockStatement()

	parser.consume(token.RIGHTBRACE, "Expect '}' after function body. got=%s", parser.peek().Type)

	return expression
}

func (parser *Parser) functionParameters() (map[string]ast.ExpressionNode, []*ast.Identifier) {
	defaults := make(map[string]ast.ExpressionNode)
	parameters := []*ast.Identifier{}

	if parser.match(token.RIGHTPAREN) {
		return defaults, parameters
	}

	// to do

	return defaults, parameters
}
