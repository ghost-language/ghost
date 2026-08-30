package object

import (
	"bytes"
	"fmt"
	"strings"

	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/token"
)

// Map objects consist of a map value.
//
// Pairs gives O(1) lookup by key, exactly like a bare Go map would; order
// tracks the sequence keys were first inserted in, separately, so a lookup
// (`get`, `has`, `[]`, `.`) can keep reading Pairs directly while every
// operation that answers more than one pair at once - keys(), values(),
// entries(), String(), `for ... in` - reads OrderedPairs() instead and
// agrees with the rest on what order that is (§13.5, §14 decision 2: Map
// guarantees insertion order, the same guarantee a JS object or a PHP
// associative array already gives their users). order is unexported so it
// can only drift out of sync with Pairs by a bug in this file - every
// mutation in the language (a literal, set(), index/property assignment,
// merge()) goes through SetPair/RemovePair, never Pairs directly.
type Map struct {
	Pairs map[MapKey]MapPair
	order []MapKey
}

type MapPair struct {
	Key   Object
	Value Object
}

// NewOrderedMap builds an empty Map ready for SetPair calls. It is the one
// way to build a Map from scratch in this codebase - a bare &Map{Pairs: ...}
// literal would leave order empty even when Pairs is not, which reads as an
// empty map to everything that iterates in order.
func NewOrderedMap() *Map {
	return &Map{Pairs: map[MapKey]MapPair{}}
}

// SetPair stores pair under hashed, preserving insertion order: an existing
// key keeps its original position and only its value changes; a new key is
// appended. This is the one door every Map mutation in the language goes
// through - see the Map doc comment.
func (mapObject *Map) SetPair(hashed MapKey, pair MapPair) {
	if _, exists := mapObject.Pairs[hashed]; !exists {
		mapObject.order = append(mapObject.order, hashed)
	}

	mapObject.Pairs[hashed] = pair
}

// RemovePair deletes hashed, answering the pair that was there and whether
// it was, and keeps order consistent for whatever remains.
func (mapObject *Map) RemovePair(hashed MapKey) (MapPair, bool) {
	pair, ok := mapObject.Pairs[hashed]

	if !ok {
		return MapPair{}, false
	}

	delete(mapObject.Pairs, hashed)

	for index, key := range mapObject.order {
		if key == hashed {
			mapObject.order = append(mapObject.order[:index], mapObject.order[index+1:]...)

			break
		}
	}

	return pair, true
}

// OrderedPairs answers every pair in the map in insertion order - the one
// place that order is read back, so keys()/values()/entries()/String() and
// `for ... in` all agree with each other and with SetPair on what it is.
func (mapObject *Map) OrderedPairs() []MapPair {
	pairs := make([]MapPair, len(mapObject.order))

	for index, key := range mapObject.order {
		pairs[index] = mapObject.Pairs[key]
	}

	return pairs
}

// String represents the map object's value as a string.
func (mapObject *Map) String() string {
	var out bytes.Buffer

	ordered := mapObject.OrderedPairs()
	pairs := make([]string, len(ordered))

	for index, pair := range ordered {
		pairs[index] = fmt.Sprintf("%s: %s", pair.Key.String(), pair.Value.String())
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

	ordered := mapObject.OrderedPairs()
	elements := make([]Object, len(ordered))

	for index, pair := range ordered {
		elements[index] = &List{Elements: []Object{pair.Key, pair.Value}}
	}

	return &List{Elements: elements}, true
}

func (mapObject *Map) keys(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("map.keys()", tok, args, 0); err != nil {
		return err, true
	}

	ordered := mapObject.OrderedPairs()
	elements := make([]Object, len(ordered))

	for index, pair := range ordered {
		elements[index] = pair.Key
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

	// This map's pairs are set first, so a key shared with other keeps this
	// map's position for it (SetPair's own rule) - the same result a plain
	// object spread (`{...this, ...other}`) gives in JS.
	merged := NewOrderedMap()

	for _, hashed := range mapObject.order {
		merged.SetPair(hashed, mapObject.Pairs[hashed])
	}

	for _, hashed := range other.order {
		merged.SetPair(hashed, other.Pairs[hashed])
	}

	return merged, true
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

	pair, ok := mapObject.RemovePair(hashed)

	if !ok {
		return &Null{}, true
	}

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

	mapObject.SetPair(hashed, MapPair{Key: key, Value: value})

	return mapObject, true
}

func (mapObject *Map) values(tok token.Token, args []Object) (Object, bool) {
	if err := Arity("map.values()", tok, args, 0); err != nil {
		return err, true
	}

	ordered := mapObject.OrderedPairs()
	elements := make([]Object, len(ordered))

	for index, pair := range ordered {
		elements[index] = pair.Value
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

// NewMap builds a Map from a plain Go map - an embedder's convenience
// constructor. Go maps carry no order of their own, so the order the result
// settles into is whatever Go's randomized map iteration gives this one
// call; it is still frozen and consistent for every read after that
// (§13.5), just not meaningfully "insertion order" the way a map literal or
// a run of set() calls is.
func NewMap(values map[string]interface{}) *Map {
	mapObject := NewOrderedMap()

	for key, value := range values {
		pairKey := &String{Value: key}
		var pairValue Object

		switch val := value.(type) {
		case int:
			pairValue = NewInt(int64(val))
		case int64:
			pairValue = NewInt(val)
		case string:
			pairValue = &String{Value: val}
		}

		mapObject.SetPair(pairKey.MapKey(), MapPair{Key: pairKey, Value: pairValue})
	}

	return mapObject
}
