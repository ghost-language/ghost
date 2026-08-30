package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) functionStatement() ast.ExpressionNode {
	expression := &ast.Function{Token: parser.currentToken}

	if !parser.nextTokenIs(token.LEFTPAREN) {
		parser.readToken()

		expression.Name = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme}
	}

	if !parser.expectNextTokenIs(token.LEFTPAREN) {
		return nil
	}

	expression.Defaults, expression.Parameters, expression.Rest = parser.functionParameters()

	if !parser.expectNextTokenIs(token.LEFTBRACE) {
		return nil
	}

	expression.Body = parser.blockStatement()

	return expression
}

// functionParameters parses a parameter list, answering the default-value
// table, the parameter names in order, and whether the last one is a rest
// parameter (`...args`, §12).
func (parser *Parser) functionParameters() (map[string]ast.ExpressionNode, []*ast.Identifier, bool) {
	defaults := make(map[string]ast.ExpressionNode)
	parameters := []*ast.Identifier{}
	rest := false

	if parser.nextTokenIs(token.RIGHTPAREN) {
		parser.readToken()

		return defaults, parameters, rest
	}

	parser.readToken()

	// As with the import loop, every turn has to consume a token and stop at the
	// end of the file: a parameter list that is never closed would otherwise
	// hang the parser rather than report itself.
	for !parser.currentTokenIs(token.RIGHTPAREN) {
		if parser.isAtEnd() {
			parser.report(parser.currentToken, "expected `)` to close the parameter list, found %s", parser.currentToken.Describe())

			return defaults, parameters, rest
		}

		if rest {
			parser.report(parser.currentToken, "a rest parameter has to be the last one")

			return defaults, parameters, rest
		}

		isRest := parser.currentTokenIs(token.DOTDOTDOT)

		if isRest {
			parser.readToken()
		}

		if !parser.currentTokenIs(token.IDENTIFIER) {
			parser.report(parser.currentToken, "expected a parameter name, found %s", parser.currentToken.Describe())

			return defaults, parameters, rest
		}

		parameter := &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Lexeme}
		parameters = append(parameters, parameter)

		parser.readToken()

		if isRest {
			rest = true

			if parser.currentTokenIs(token.EQUAL) {
				parser.report(parser.currentToken, "a rest parameter cannot have a default value")

				return defaults, parameters, rest
			}
		} else if parser.currentTokenIs(token.EQUAL) {
			parser.readToken()

			defaults[parameter.Value] = parser.expressionStatement()

			parser.readToken()
		}

		if parser.currentTokenIs(token.COMMA) {
			parser.readToken()
		}
	}

	return defaults, parameters, rest
}
