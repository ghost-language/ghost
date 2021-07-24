package object

type Callable interface {
	Arity() int
	Call([]interface{}) (interface{}, error)
}