package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

// destructuringAssign parses a statement beginning with `[` or `{`. The
// grammar for a pattern (`[a, b]`, `{x, y}`, `{x: a}`) is exactly the grammar
// for a list or map literal - parsing doesn't evaluate anything, so there is
// no harm reading it as one first. What follows decides which it actually
// was: `=` makes it a destructuring assignment, anything else means it was a
// literal expression all along (`[1, 2, 3]` as a statement on its own, say).
func (parser *Parser) destructuringAssign() ast.StatementNode {
	expression := parser.parseExpression(LOWEST)

	if !parser.nextTokenIs(token.EQUAL) {
		return &ast.Expression{Expression: expression}
	}

	var name ast.AssignmentNode

	switch pattern := expression.(type) {
	case *ast.List:
		targets, ok := parser.listPatternTargets(pattern)

		if !ok {
			return &ast.Expression{Expression: expression}
		}

		name = &ast.ListPattern{Token: pattern.Token, Targets: targets}
	case *ast.Map:
		pairs, ok := parser.mapPatternPairs(pattern)

		if !ok {
			return &ast.Expression{Expression: expression}
		}

		name = &ast.MapPattern{Token: pattern.Token, Pairs: pairs}
	default:
		return &ast.Expression{Expression: expression}
	}

	statement := &ast.Assign{Name: name}

	parser.readToken() // consume up to `=`
	statement.Token = parser.currentToken
	parser.readToken()

	statement.Value = parser.parseExpression(LOWEST)

	if parser.nextTokenIs(token.SEMICOLON) {
		parser.readToken()
	}

	return statement
}

// listPatternTargets requires every element of an already-parsed list
// literal to be a plain name - `[a, b+1] = list` names something that isn't
// a valid assignment target, the same restriction plain assignment already
// puts on its left-hand side.
func (parser *Parser) listPatternTargets(list *ast.List) ([]*ast.Identifier, bool) {
	targets := make([]*ast.Identifier, 0, len(list.Elements))

	for _, element := range list.Elements {
		identifier, ok := element.(*ast.Identifier)

		if !ok {
			parser.report(list.Token, "a list pattern can only bind plain names, not a full expression")

			return nil, false
		}

		targets = append(targets, identifier)
	}

	return targets, true
}

// mapPatternPairs requires every pair of an already-parsed map literal to
// read a plain key into a plain name - the shorthand `{x, y}` already parses
// this way (map.go's shorthand-key handling sets the value to a copy of the
// key identifier), and `{x: a}` renames the binding the same way an explicit
// map literal pair can name any value.
func (parser *Parser) mapPatternPairs(mapLiteral *ast.Map) ([]ast.MapPatternPair, bool) {
	pairs := make([]ast.MapPatternPair, 0, len(mapLiteral.Pairs))

	for key, value := range mapLiteral.Pairs {
		source, ok := key.(*ast.Identifier)

		if !ok {
			parser.report(mapLiteral.Token, "a map pattern can only bind from a plain name, not a computed key")

			return nil, false
		}

		target, ok := value.(*ast.Identifier)

		if !ok {
			parser.report(mapLiteral.Token, "a map pattern can only bind to a plain name, not a full expression")

			return nil, false
		}

		pairs = append(pairs, ast.MapPatternPair{Source: source, Target: target})
	}

	return pairs, true
}
