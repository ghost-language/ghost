package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) assign() ast.StatementNode {
	if parser.check(token.IDENTIFIER) && parser.checkNext(token.ASSIGN) {
		statement := &ast.Assign{}

		statement.Name = &ast.Identifier{Token: parser.peek(), Value: parser.peek().Lexeme}
		parser.match(token.IDENTIFIER)
		statement.Token = parser.peek()
		parser.match(token.ASSIGN)

		value := parser.parseExpression(LOWEST)

		if value == nil {
			return nil
		}

		statement.Value = value

		return statement
	}

	return nil
}
