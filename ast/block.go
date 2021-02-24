package ast

type Block struct {
	StatementNode
	Statements []StatementNode
}
