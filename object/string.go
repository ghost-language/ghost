package object

import "hash/fnv"

type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}

func (s *String) Inspect() string {
	return s.Value
}

func (s *String) Set(obj Object) {
	s.Value = obj.(*String).Value
}

// MapKey defines the key for maps that can be comparable and unique.
//
// Note: There is a _very_ small chance that the following will
// result in the same hash being generated for different string
// values (hash collisions). Research "separate chaining" and
// "open addressing" techniques to work around the problem.
func (s *String) MapKey() MapKey {
	hash := fnv.New64a()
	hash.Write([]byte(s.Value))

	return MapKey{Type: s.Type(), Value: hash.Sum64()}
}
