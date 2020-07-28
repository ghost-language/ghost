package object

import "fmt"

type BuiltinFunction func(env *Environment, args ...Object) Object

type Builtin struct {
	Name string
	Fn   BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}

func (b *Builtin) Inspect() string {
	return fmt.Sprintf("builtin function: %s", b.Name)
}
