package parser

import (
	"fmt"
	"strconv"

	"ghostlang.org/ghost/ast"
)

func (p *Parser) parseIntegerLiteral() ast.Expression {
	integerLiteral := &ast.IntegerLiteral{Token: p.currentToken}

	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)

	if err != nil {
		message := fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, message)

		return nil
	}

	integerLiteral.Value = value

	return integerLiteral
}
