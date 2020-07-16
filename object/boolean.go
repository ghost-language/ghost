package object

import "fmt"

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) Set(obj Object) {
	b.Value = obj.(*Boolean).Value
}

func (b *Boolean) MapKey() MapKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return MapKey{Type: b.Type(), Value: value}
}
