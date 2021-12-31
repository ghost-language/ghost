package ast

type Node interface {
	Visitable
}

type StatementNode interface {
	Node
}

type ExpressionNode interface {
	Node
}

type Visitor interface {
	visitAssign(*Assign)
	visitBlock(*Block)
	visitBoolean(*Boolean)
	visitCall(*Call)
	visitExpression(*Expression)
	visitFunction(*Function)
	visitIdentifier(*Identifier)
	visitIf(*If)
	visitIndex(*Index)
	visitInfix(*Infix)
	visitList(*List)
	visitMap(*Map)
	visitMethod(*Method)
	visitNull(*Null)
	visitNumber(*Number)
	visitPrefix(*Prefix)
	visitProgram(*Program)
	visitString(*String)
	visitWhile(*While)
}

type Visitable interface {
	Accept(Visitor)
}
