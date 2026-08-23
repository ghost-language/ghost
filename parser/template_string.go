package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

// templateLiteral builds a TemplateString from the chunk/expression token
// stream the scanner produces for a backtick literal: a TEMPLATESTRING chunk
// for each piece of text followed by an interpolation, ending in a single
// TEMPLATESTRINGEND chunk for the text that closes the literal.
func (parser *Parser) templateLiteral() ast.ExpressionNode {
	template := &ast.TemplateString{Token: parser.currentToken}
	template.Chunks = append(template.Chunks, parser.currentToken.Literal.(string))

	for parser.currentTokenIs(token.TEMPLATESTRING) {
		parser.readToken()

		template.Expressions = append(template.Expressions, parser.parseExpression(LOWEST))

		if !parser.nextTokenIs(token.TEMPLATESTRING) && !parser.nextTokenIs(token.TEMPLATESTRINGEND) {
			parser.report(parser.nextToken, "expected the template literal to continue, found %s", parser.nextToken.Describe())

			return template
		}

		parser.readToken()

		template.Chunks = append(template.Chunks, parser.currentToken.Literal.(string))
	}

	return template
}
