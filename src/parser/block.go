package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) blockStatement() *ast.Block {
	block := &ast.Block{Token: parser.currentToken}
	block.Statements = []ast.StatementNode{}

	parser.readToken()

	for !parser.currentTokenIs(token.RIGHTBRACE) && !parser.isAtEnd() {
		statement := parser.statement()

		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}

		parser.readToken()
	}

	return block
}
