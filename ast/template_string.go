package ast

import (
	"ghostlang.org/x/ghost/token"
)

// TemplateString represents a backtick-delimited template literal: literal
// text chunks interleaved with expressions to interpolate, e.g.
// `count: ${count}`. Chunks always has one more entry than Expressions — the
// text before the first interpolation, the text between each pair of
// interpolations, and the text after the last one.
type TemplateString struct {
	ExpressionNode
	Token       token.Token
	Chunks      []string
	Expressions []ExpressionNode
}
