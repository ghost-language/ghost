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
	case "charAt":
		return str.charAt(tok, args)
	case "contains":
		return str.contains(tok, args)
	case "endsWith":
		return str.endsWith(tok, args)
	case "find":
		return str.find(tok, args)
	case "findAll":
		return str.findAll(tok, args)
	case "format":
		return str.format(tok, args)
	case "indexOf":
		return str.indexOf(tok, args)
	case "isEmpty":
		return str.isEmpty(tok, args)
	case "lastIndexOf":
		return str.lastIndexOf(tok, args)
	case "length":
		return str.length(tok, args)
	case "matches":
		return str.matches(tok, args)
	case "padEnd":
		return str.padEnd(tok, args)
	case "padStart":
		return str.padStart(tok, args)
	case "repeat":
		return str.repeat(tok, args)
	case "replace":
		return str.replace(tok, args)
	case "reverse":
		return str.reverse(tok, args)
	case "slice":
		return str.slice(tok, args)
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

// runes is the rune-indexed view of the string, so a method that names a
// position or a range agrees with length() about what a "character" is -
// one multi-byte rune, not one byte.
func (str *String) runes() []rune {
	return []rune(str.Value)
}

// charAt answers the single-character string at a rune position, or an empty
// string for a position out of range - the same leniency list indexing
// already gives a position that names a spot rather than a range (§13.6).
func (str *String) charAt(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.charAt()", tok, args, 1); err != nil {
		return err, true
	}

	index, err := NumberArgument("string.charAt()", tok, args, 0)

	if err != nil {
		return err, true
	}

	runes := str.runes()
	idx := index.Int64()

	if idx < 0 || idx >= int64(len(runes)) {
		return &String{}, true
	}

	return &String{Value: string(runes[idx])}, true
}

// contains reports whether a substring appears anywhere in the string,
// mirroring list.contains() so "does this collection have this value" reads
// the same on both types.
func (str *String) contains(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.contains()", tok, args, 1); err != nil {
		return err, true
	}

	substr, err := StringArgument("string.contains()", tok, args, 0)

	if err != nil {
		return err, true
	}

	return &Boolean{Value: strings.Contains(str.Value, substr.Value)}, true
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

// indexOf answers the rune position of the first occurrence of a substring,
// or -1 if it never appears - the byte offset strings.Index finds is
// converted to a rune count so it lines up with charAt() and length().
func (str *String) indexOf(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.indexOf()", tok, args, 1); err != nil {
		return err, true
	}

	substr, err := StringArgument("string.indexOf()", tok, args, 0)

	if err != nil {
		return err, true
	}

	byteIndex := strings.Index(str.Value, substr.Value)

	if byteIndex < 0 {
		return NewInt(-1), true
	}

	return NewInt(int64(utf8.RuneCountInString(str.Value[:byteIndex]))), true
}

func (str *String) isEmpty(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.isEmpty()", tok, args, 0); err != nil {
		return err, true
	}

	return &Boolean{Value: len(str.Value) == 0}, true
}

// lastIndexOf answers the rune position of the last occurrence of a
// substring, or -1 if it never appears.
func (str *String) lastIndexOf(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.lastIndexOf()", tok, args, 1); err != nil {
		return err, true
	}

	substr, err := StringArgument("string.lastIndexOf()", tok, args, 0)

	if err != nil {
		return err, true
	}

	byteIndex := strings.LastIndex(str.Value, substr.Value)

	if byteIndex < 0 {
		return NewInt(-1), true
	}

	return NewInt(int64(utf8.RuneCountInString(str.Value[:byteIndex]))), true
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

// pad grows a string to a target rune length by adding copies of a pad
// string to one side, truncated to fit exactly. A string already at or past
// the target length, or an empty pad string, comes back unchanged - there is
// nothing to repeat that would make it longer.
func pad(runes []rune, targetLength int64, padStr string, before bool) string {
	deficit := targetLength - int64(len(runes))

	if deficit <= 0 || padStr == "" {
		return string(runes)
	}

	padRunes := []rune(padStr)
	fill := make([]rune, 0, deficit)

	for int64(len(fill)) < deficit {
		fill = append(fill, padRunes...)
	}

	fill = fill[:deficit]

	if before {
		return string(fill) + string(runes)
	}

	return string(runes) + string(fill)
}

func (str *String) padEnd(tok token.Token, args []Object) (Object, bool) {
	if err := ArityRange("string.padEnd()", tok, args, 1, 2); err != nil {
		return err, true
	}

	targetLength, err := NumberArgument("string.padEnd()", tok, args, 0)

	if err != nil {
		return err, true
	}

	padStr := " "

	if len(args) == 2 {
		padArg, err := StringArgument("string.padEnd()", tok, args, 1)

		if err != nil {
			return err, true
		}

		padStr = padArg.Value
	}

	return &String{Value: pad(str.runes(), targetLength.Int64(), padStr, false)}, true
}

func (str *String) padStart(tok token.Token, args []Object) (Object, bool) {
	if err := ArityRange("string.padStart()", tok, args, 1, 2); err != nil {
		return err, true
	}

	targetLength, err := NumberArgument("string.padStart()", tok, args, 0)

	if err != nil {
		return err, true
	}

	padStr := " "

	if len(args) == 2 {
		padArg, err := StringArgument("string.padStart()", tok, args, 1)

		if err != nil {
			return err, true
		}

		padStr = padArg.Value
	}

	return &String{Value: pad(str.runes(), targetLength.Int64(), padStr, true)}, true
}

func (str *String) repeat(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.repeat()", tok, args, 1); err != nil {
		return err, true
	}

	count, err := NumberArgument("string.repeat()", tok, args, 0)

	if err != nil {
		return err, true
	}

	times := count.Int64()

	if times < 0 {
		return NewError(fault.Value, tok, "`string.repeat()` count cannot be negative, got %d", times), true
	}

	return &String{Value: strings.Repeat(str.Value, int(times))}, true
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

func (str *String) reverse(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("string.reverse()", tok, args, 0); err != nil {
		return err, true
	}

	runes := str.runes()
	length := len(runes)
	reversed := make([]rune, length)

	for index, r := range runes {
		reversed[length-1-index] = r
	}

	return &String{Value: string(reversed)}, true
}

// slice answers a new string holding the runes from start up to, but not
// including, end - which defaults to the length of the string. Bounds are
// validated the same way list.slice() validates them (§13.6): a range names
// two positions, so both are checked rather than clamped.
func (str *String) slice(tok token.Token, args []Object) (Object, bool) {
	if err := ArityRange("string.slice()", tok, args, 1, 2); err != nil {
		return err, true
	}

	start, err := NumberArgument("string.slice()", tok, args, 0)

	if err != nil {
		return err, true
	}

	runes := str.runes()
	length := int64(len(runes))
	from := start.Int64()
	to := length

	if len(args) == 2 {
		end, err := NumberArgument("string.slice()", tok, args, 1)

		if err != nil {
			return err, true
		}

		to = end.Int64()
	}

	if from < 0 || from > length {
		return NewError(fault.Index, tok, "`string.slice()` start index %d is out of range for a string of length %d", from, length), true
	}

	if to < from || to > length {
		return NewError(fault.Index, tok, "`string.slice()` end index %d is out of range for a string of length %d", to, length), true
	}

	return &String{Value: string(runes[from:to])}, true
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
