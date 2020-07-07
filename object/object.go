package object

type ObjectType string

const (
	BOOLEAN_OBJ      = "BOOLEAN"
	INTEGER_OBJ      = "INTEGER"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}
