package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) blockStatement() *ast.Block {
	block := &ast.Block{Token: parser.current()}
	block.Statements = []ast.StatementNode{}

	parser.advance()

	for !parser.check(token.RIGHTBRACE) && !parser.isAtEnd() {
		statement := parser.statement()

		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}

		parser.advance()
	}

	return block
}
