package object

import (
	"bytes"
	"fmt"
	"strings"

	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/token"
)

// Map objects consist of a map value.
type Map struct {
	Pairs map[MapKey]MapPair
}

type MapPair struct {
	Key   Object
	Value Object
}

// String represents the map object's value as a string.
func (mapObject *Map) String() string {
	var out bytes.Buffer

	length := len(mapObject.Pairs)
	pairs := make([]string, length)

	var index int

	for _, pair := range mapObject.Pairs {
		pairs[index] = fmt.Sprintf("%s: %s", pair.Key.String(), pair.Value.String())
		index++
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// Type returns the map object type.
func (mapObject *Map) Type() Type {
	return MAP
}

// Method defines the set of methods available on map objects.
func (mapObject *Map) Method(method string, tok token.Token, args []Object) (Object, bool) {
	switch method {
	case "entries":
		return mapObject.entries(tok, args)
	case "get":
		return mapObject.get(tok, args)
	case "has":
		return mapObject.has(tok, args)
	case "keys":
		return mapObject.keys(tok, args)
	case "length":
		return mapObject.length(tok, args)
	case "merge":
		return mapObject.merge(tok, args)
	case "remove":
		return mapObject.remove(tok, args)
	case "set":
		return mapObject.set(tok, args)
	case "values":
		return mapObject.values(tok, args)
	}

	return nil, false
}

// =============================================================================
// Object methods

// get reads the value stored under a key, or a default when the key is
// absent - and null when there is no default either, the same as indexing
// with `[]` answers for a missing key.
func (mapObject *Map) get(tok token.Token, args []Object) (Object, bool) {
	if err := ArityRange("map.get()", tok, args, 1, 2); err != nil {
		return err, true
	}

	_, hashed, err := mapKey("map.get()", tok, args, 0)

	if err != nil {
		return err, true
	}

	if pair, ok := mapObject.Pairs[hashed]; ok {
		return pair.Value, true
	}

	if len(args) == 2 {
		return args[1], true
	}

	return &Null{}, true
}

// has reports whether a key is present, regardless of what its value is -
// which is the distinction get() alone cannot make, since a key can map to
// null on purpose.
func (mapObject *Map) has(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("map.has()", tok, args, 1); err != nil {
		return err, true
	}

	_, hashed, err := mapKey("map.has()", tok, args, 0)

	if err != nil {
		return err, true
	}

	_, ok := mapObject.Pairs[hashed]

	return &Boolean{Value: ok}, true
}

// entries answers a list of [key, value] pairs, one per entry - the same
// shape keys()/values() would zip together, kept here as a single call.
func (mapObject *Map) entries(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("map.entries()", tok, args, 0); err != nil {
		return err, true
	}

	elements := make([]Object, 0, len(mapObject.Pairs))

	for _, pair := range mapObject.Pairs {
		elements = append(elements, &List{Elements: []Object{pair.Key, pair.Value}})
	}

	return &List{Elements: elements}, true
}

func (mapObject *Map) keys(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("map.keys()", tok, args, 0); err != nil {
		return err, true
	}

	elements := make([]Object, 0, len(mapObject.Pairs))

	for _, pair := range mapObject.Pairs {
		elements = append(elements, pair.Key)
	}

	return &List{Elements: elements}, true
}

func (mapObject *Map) length(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("map.length()", tok, args, 0); err != nil {
		return err, true
	}

	return NewInt(int64(len(mapObject.Pairs))), true
}

// merge answers a new map holding this map's pairs and another's, keeping
// this map untouched. Where a key appears in both, the other map's value
// wins - the same rule a later assignment to the same key would follow.
func (mapObject *Map) merge(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("map.merge()", tok, args, 1); err != nil {
		return err, true
	}

	other, err := MapArgument("map.merge()", tok, args, 0)

	if err != nil {
		return err, true
	}

	pairs := make(map[MapKey]MapPair, len(mapObject.Pairs)+len(other.Pairs))

	for key, pair := range mapObject.Pairs {
		pairs[key] = pair
	}

	for key, pair := range other.Pairs {
		pairs[key] = pair
	}

	return &Map{Pairs: pairs}, true
}

// remove deletes a key from the map, mutating in place, and answers the
// value that was stored there - null if the key was not present, the same
// leniency pop()/shift() give an empty list rather than an error.
func (mapObject *Map) remove(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("map.remove()", tok, args, 1); err != nil {
		return err, true
	}

	_, hashed, err := mapKey("map.remove()", tok, args, 0)

	if err != nil {
		return err, true
	}

	pair, ok := mapObject.Pairs[hashed]

	if !ok {
		return &Null{}, true
	}

	delete(mapObject.Pairs, hashed)

	return pair.Value, true
}

// set stores a value under a key, adding it if the key is new, and answers
// the map itself so a call can chain.
func (mapObject *Map) set(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("map.set()", tok, args, 2); err != nil {
		return err, true
	}

	key, hashed, err := mapKey("map.set()", tok, args, 0)

	if err != nil {
		return err, true
	}

	value, valueErr := AnyArgument("map.set()", tok, args, 1)

	if valueErr != nil {
		return valueErr, true
	}

	mapObject.Pairs[hashed] = MapPair{Key: key, Value: value}

	return mapObject, true
}

func (mapObject *Map) values(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("map.values()", tok, args, 0); err != nil {
		return err, true
	}

	elements := make([]Object, 0, len(mapObject.Pairs))

	for _, pair := range mapObject.Pairs {
		elements = append(elements, pair.Value)
	}

	return &List{Elements: elements}, true
}

// mapKey reads an argument that has to work as a map key, which is a
// narrower requirement than being any particular type - the same one
// indexing a map with `[]` enforces. It answers the key hashed and ready to
// use, alongside the object itself for storing back in a MapPair.
func mapKey(name string, tok token.Token, args []Object, index int) (Object, MapKey, *Error) {
	value, err := AnyArgument(name, tok, args, index)

	if err != nil {
		return nil, MapKey{}, err
	}

	mappable, ok := value.(Mappable)

	if !ok {
		return nil, MapKey{}, NewError(fault.Type, tok, "`%s` cannot use %s as a map key", name, TypeName(value)).
			WithHelp("a map key has to be a string, a number, or a boolean")
	}

	return value, mappable.MapKey(), nil
}

func NewMap(values map[string]interface{}) *Map {
	pairs := make(map[MapKey]MapPair)

	for key, value := range values {
		pairKey := &String{Value: key}
		var pairValue Object
		hashed := pairKey.MapKey()

		switch val := value.(type) {
		case int:
			pairValue = NewInt(int64(val))
		case int64:
			pairValue = NewInt(val)
		case string:
			pairValue = &String{Value: val}
		}

		pairs[hashed] = MapPair{Key: pairKey, Value: pairValue}
	}

	return &Map{Pairs: pairs}
}
