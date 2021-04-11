package object

import "fmt"

type NativeFunction struct {
	Callable
	nativeCall GhostCallable
	arity int
}

func (n *NativeFunction) Call(arguments []Object) (Object, error) {
	return n.nativeCall(arguments)
}

func (n *NativeFunction) Arity() int {
	return n.arity
}

func (n *NativeFunction) String() string {
	return fmt.Sprintf("native function: %p", n.nativeCall)
}