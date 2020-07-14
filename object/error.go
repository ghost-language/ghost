package object

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

func (e *Error) Set(obj Object) {
	e.Message = obj.(*Error).Message
}
