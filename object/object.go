package object

type ObjectType string

const (
	ARRAY_OBJ        = "ARRAY"
	BOOLEAN_OBJ      = "BOOLEAN"
	BUILTIN_OBJ      = "BUILTIN"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	NULL_OBJ         = "NULL"
	NUMBER_OBJ       = "NUMBER"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	STRING_OBJ       = "STRING"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}
