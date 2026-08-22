package object

import (
	"fmt"
	"hash/fnv"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/token"
)

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
//
// Every method here reads its arguments through the shared helpers rather than
// asserting types inline. A method called with the wrong arguments answers with
// an argument error pointing at the call, which is the whole reason the token
// is threaded this far down.
func (str *String) Method(method string, tok token.Token, args []Object) (Object, bool) {
	switch method {
	case "find":
		return str.find(tok, args)
	case "findAll":
		return str.findAll(tok, args)
	case "format":
		return str.format(tok, args)
	case "endsWith":
		return str.endsWith(tok, args)
	case "length":
		return str.length(tok, args)
	case "matches":
		return str.matches(tok, args)
	case "replace":
		return str.replace(tok, args)
	case "split":
		return str.split(tok, args)
	case "startsWith":
		return str.startsWith(tok, args)
	case "toLowerCase":
		return str.toLowerCase(tok, args)
	case "toUpperCase":
		return str.toUpperCase(tok, args)
	case "toString":
		return str.toString(tok, args)
	case "toNumber":
		return str.toNumber(tok, args)
	case "trim":
		return str.trim(tok, args)
	case "trimEnd":
		return str.trimEnd(tok, args)
	case "trimStart":
		return str.trimStart(tok, args)
	}

	return nil, false
}

// =============================================================================
// Object methods

// pattern compiles the string as a regular expression. A string is only a
// pattern when a method treats it as one, so a bad pattern is reported at the
// call that made that assumption rather than crashing the interpreter.
func (str *String) pattern(name string, tok token.Token) (*regexp.Regexp, *Error) {
	compiled, err := regexp.Compile(str.Value)

	if err != nil {
		return nil, NewError(fault.Value, tok, "`%s` cannot use `%s` as a pattern: %s", name, str.Value, patternReason(err)).
			WithHelp("the string a pattern method is called on is the pattern itself")
	}

	return compiled, nil
}

// patternReason strips the Go wrapping from a regular expression error, leaving
// the part that describes what is wrong with the pattern. Go's message opens
// with a label and closes by quoting the pattern back; the message being built
// around it already has both.
func patternReason(err error) string {
	reason := strings.TrimPrefix(err.Error(), "error parsing regexp: ")

	if index := strings.LastIndex(reason, ": `"); index >= 0 {
		return reason[:index]
	}

	return reason
}

func (str *String) find(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.find()", tok, args, 1); err != nil {
		return err, true
	}

	subject, err := StringArgument("string.find()", tok, args, 0)

	if err != nil {
		return err, true
	}

	expression, err := str.pattern("string.find()", tok)

	if err != nil {
		return err, true
	}

	found := expression.FindStringSubmatch(subject.Value)

	if len(found) > 0 {
		return &String{Value: found[0]}, true
	}

	return &String{}, true
}

func (str *String) findAll(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.findAll()", tok, args, 1); err != nil {
		return err, true
	}

	subject, err := StringArgument("string.findAll()", tok, args, 0)

	if err != nil {
		return err, true
	}

	expression, err := str.pattern("string.findAll()", tok)

	if err != nil {
		return err, true
	}

	list := &List{}

	for _, match := range expression.FindStringSubmatch(subject.Value) {
		list.Elements = append(list.Elements, &String{Value: match})
	}

	return list, true
}

func (str *String) format(tok token.Token, args []Object) (Object, bool) {
	values := make([]interface{}, 0, len(args))

	for _, value := range args {
		values = append(values, value.String())
	}

	return &String{Value: fmt.Sprintf(str.Value, values...)}, true
}

func (str *String) endsWith(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.endsWith()", tok, args, 1); err != nil {
		return err, true
	}

	suffix, err := StringArgument("string.endsWith()", tok, args, 0)

	if err != nil {
		return err, true
	}

	return &Boolean{Value: strings.HasSuffix(str.Value, suffix.Value)}, true
}

func (str *String) length(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.length()", tok, args, 0); err != nil {
		return err, true
	}

	return NewInt(int64(utf8.RuneCountInString(str.Value))), true
}

func (str *String) matches(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.matches()", tok, args, 1); err != nil {
		return err, true
	}

	subject, err := StringArgument("string.matches()", tok, args, 0)

	if err != nil {
		return err, true
	}

	expression, err := str.pattern("string.matches()", tok)

	if err != nil {
		return err, true
	}

	return &Boolean{Value: expression.MatchString(subject.Value)}, true
}

func (str *String) replace(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.replace()", tok, args, 2); err != nil {
		return err, true
	}

	from, err := StringArgument("string.replace()", tok, args, 0)

	if err != nil {
		return err, true
	}

	to, err := StringArgument("string.replace()", tok, args, 1)

	if err != nil {
		return err, true
	}

	return &String{Value: strings.ReplaceAll(str.Value, from.Value, to.Value)}, true
}

func (str *String) split(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.split()", tok, args, 1); err != nil {
		return err, true
	}

	separator, err := StringArgument("string.split()", tok, args, 0)

	if err != nil {
		return err, true
	}

	parts := strings.Split(str.Value, separator.Value)
	elements := make([]Object, len(parts))

	for index, part := range parts {
		elements[index] = &String{Value: part}
	}

	return &List{Elements: elements}, true
}

func (str *String) startsWith(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.startsWith()", tok, args, 1); err != nil {
		return err, true
	}

	prefix, err := StringArgument("string.startsWith()", tok, args, 0)

	if err != nil {
		return err, true
	}

	return &Boolean{Value: strings.HasPrefix(str.Value, prefix.Value)}, true
}

func (str *String) toLowerCase(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.toLowerCase()", tok, args, 0); err != nil {
		return err, true
	}

	return &String{Value: strings.ToLower(str.Value)}, true
}

func (str *String) toUpperCase(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.toUpperCase()", tok, args, 0); err != nil {
		return err, true
	}

	return &String{Value: strings.ToUpper(str.Value)}, true
}

func (str *String) toString(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.toString()", tok, args, 0); err != nil {
		return err, true
	}

	return str, true
}

func (str *String) toNumber(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.toNumber()", tok, args, 0); err != nil {
		return err, true
	}

	trimmed := strings.TrimSpace(str.Value)

	if whole, err := strconv.ParseInt(trimmed, 10, 64); err == nil {
		return NewInt(whole), true
	}

	fractional, err := strconv.ParseFloat(trimmed, 64)

	if err != nil {
		return NewError(fault.Value, tok, "`string.toNumber()` cannot read `%s` as a number", str.Value).
			WithHelp("a number is digits, optionally with a sign, a decimal point, or an exponent"), true
	}

	return NewFloat(fractional), true
}

func (str *String) trim(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.trim()", tok, args, 0); err != nil {
		return err, true
	}

	return &String{Value: strings.TrimSpace(str.Value)}, true
}

func (str *String) trimEnd(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.trimEnd()", tok, args, 0); err != nil {
		return err, true
	}

	return &String{Value: strings.TrimRight(str.Value, "\t\n\v\f\r ")}, true
}

func (str *String) trimStart(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.trimStart()", tok, args, 0); err != nil {
		return err, true
	}

	return &String{Value: strings.TrimLeft(str.Value, "\t\n\v\f\r ")}, true
}
