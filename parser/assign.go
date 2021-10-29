package parser

import (
	"fmt"

	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) assign() ast.ExpressionNode {
	if parser.check(token.IDENTIFIER) && parser.checkNext(token.ASSIGN) {
		parser.match(token.IDENTIFIER)
		name := parser.previous()
		parser.match(token.ASSIGN)

		value := parser.expression()

		return &ast.Assign{Name: name, Value: value}
	}

	log.LogDebug("failed!", fmt.Sprintf("%v -> %v : %v", parser.peek(), parser.next(), parser.current))

	return nil
}
