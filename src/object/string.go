package object

import "hash/fnv"

const STRING = "STRING"

// String objects consist of a string value.
type String struct {
	Value string
}

func (object *String) Accept(v Visitor) {
	v.visitString(object)
}

// String represents the string object's value as a string. So meta.
func (string *String) String() string {
	return string.Value
}

// Type returns the string object type.
func (string *String) Type() Type {
	return STRING
}

// MapKey defines a unique hash value for use as a map key.
func (string *String) MapKey() MapKey {
	hash := fnv.New64a()

	hash.Write([]byte(string.Value))

	return MapKey{Type: string.Type(), Value: hash.Sum64()}
}
