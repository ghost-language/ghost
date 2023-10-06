package object

import (
	"fmt"
	"hash/fnv"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/shopspring/decimal"
)

const STRING = "STRING"

// String objects consist of a string value.
type String struct {
	Value string
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

// Method defines the set of methods available on string objects.
func (str *String) Method(method string, args []Object) (Object, bool) {
	switch method {
	case "find":
		return str.find(args)
	case "findAll":
		return str.findAll(args)
	case "format":
		return str.format(args)
	case "endsWith":
		return str.endsWith(args)
	case "length":
		return str.length(args)
	case "matches":
		return str.matches(args)
	case "replace":
		return str.replace(args)
	case "split":
		return str.split(args)
	case "startsWith":
		return str.startsWith(args)
	case "toLowerCase":
		return str.toLowerCase(args)
	case "toUpperCase":
		return str.toUpperCase(args)
	case "toString":
		return str.toString(args)
	case "toNumber":
		return str.toNumber(args)
	case "trim":
		return str.trim(args)
	case "trimEnd":
		return str.trimEnd(args)
	case "trimStart":
		return str.trimStart(args)
	}

	return nil, false
}

// =============================================================================
// Object methods

func (str *String) find(args []Object) (Object, bool) {
	re := regexp.MustCompile(str.Value)

	found := re.FindStringSubmatch(args[0].(*String).Value)

	if len(found) > 0 {
		return &String{Value: found[1]}, true
	}

	return &String{}, true
}

func (str *String) findAll(args []Object) (Object, bool) {
	re := regexp.MustCompile(str.Value)
	list := &List{}
	found := re.FindStringSubmatch(args[0].(*String).Value)

	for _, f := range found {
		list.Elements = append(list.Elements, &String{Value: f})
	}

	return list, true
}

func (str *String) format(args []Object) (Object, bool) {
	list := []interface{}{}

	for _, value := range args {
		list = append(list, value.String())
	}

	return &String{Value: fmt.Sprintf(str.Value, list...)}, true
}

func (str *String) endsWith(args []Object) (Object, bool) {
	hasSuffix := strings.HasSuffix(str.Value, args[0].(*String).Value)

	return &Boolean{Value: hasSuffix}, true
}

func (str *String) length(args []Object) (Object, bool) {
	length := &Number{Value: decimal.NewFromInt(int64(utf8.RuneCountInString(str.Value)))}

	return length, true
}

func (str *String) matches(args []Object) (Object, bool) {
	matches, err := regexp.Match(str.Value, []byte(args[0].(*String).Value))

	if err != nil {
		return &Error{Message: err.Error()}, false
	}

	return &Boolean{Value: matches}, true
}

func (str *String) replace(args []Object) (Object, bool) {
	value := strings.Replace(str.Value, args[0].(*String).Value, args[1].(*String).Value, -1)

	return &String{Value: value}, true
}

func (str *String) split(args []Object) (Object, bool) {
	split := strings.Split(str.Value, args[0].(*String).Value)
	list := &List{}

	for _, value := range split {
		list.Elements = append(list.Elements, &String{Value: value})
	}

	return list, true
}

func (str *String) startsWith(args []Object) (Object, bool) {
	hasPrefix := strings.HasPrefix(str.Value, args[0].(*String).Value)

	return &Boolean{Value: hasPrefix}, true
}

func (str *String) toLowerCase(args []Object) (Object, bool) {
	return &String{Value: strings.ToLower(str.Value)}, true
}

func (str *String) toUpperCase(args []Object) (Object, bool) {
	return &String{Value: strings.ToUpper(str.Value)}, true
}

func (str *String) toString(args []Object) (Object, bool) {
	return str, true
}

func (str *String) toNumber(args []Object) (Object, bool) {
	number, _ := decimal.NewFromString(str.Value)

	return &Number{Value: number}, true
}

func (str *String) trim(args []Object) (Object, bool) {
	return &String{Value: strings.TrimSpace(str.Value)}, true
}

func (str *String) trimEnd(args []Object) (Object, bool) {
	return &String{Value: strings.TrimRight(str.Value, "\t\n\v\f\r ")}, true
}

func (str *String) trimStart(args []Object) (Object, bool) {
	return &String{Value: strings.TrimLeft(str.Value, "\t\n\v\f\r ")}, true
}
