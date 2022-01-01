package object

import (
	"hash/fnv"
	"unicode/utf8"

	"github.com/shopspring/decimal"
)

const STRING = "STRING"

// String objects consist of a string value.
type String struct {
	Value string
}

func (str *String) Accept(v Visitor) {
	v.visitString(str)
}

// String represents the string object's value as a string. So meta.
func (str *String) String() string {
	return str.Value
}

// Type returns the string object type.
func (str *String) Type() Type {
	return STRING
}

// MapKey defines a unique hash value for use as a map key.
func (str *String) MapKey() MapKey {
	hash := fnv.New64a()

	hash.Write([]byte(str.Value))

	return MapKey{Type: str.Type(), Value: hash.Sum64()}
}

func (str *String) Method(method string, args []Object) (Object, bool) {
	switch method {
	case "length":
		return str.length(args)
	}

	return nil, false
}

// Methods

func (str *String) length(args []Object) (Object, bool) {
	length := &Number{Value: decimal.NewFromInt(int64(utf8.RuneCountInString(str.Value)))}

	return length, true
}
