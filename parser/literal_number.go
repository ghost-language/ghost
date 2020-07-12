package parser

import (
	"fmt"

	"ghostlang.org/ghost/ast"
	"ghostlang.org/ghost/decimal"
)

func (p *Parser) parseNumberLiteral() ast.Expression {
	numberLiteral := &ast.NumberLiteral{Token: p.currentToken}

	value, err := decimal.NewFromString(p.currentToken.Literal)
	if err != nil {
		message := fmt.Sprintf("could not parse %q as number", p.currentToken.Literal)
		p.errors = append(p.errors, message)

		return nil
	}

	numberLiteral.Value = value

	return numberLiteral
}
