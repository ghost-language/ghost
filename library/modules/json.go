package modules

import (
	"encoding/json"
	"fmt"

	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

var JsonMethods = map[string]*object.LibraryFunction{}
var JsonProperties = map[string]*object.LibraryProperty{}

func init() {
	RegisterMethod(JsonMethods, "decode", jsonDecode)
	RegisterMethod(JsonMethods, "encode", jsonEncode)
}

// jsonDecode reads JSON text into a list or a map.
func jsonDecode(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("json.decode", tok, args, 1); err != nil {
		return err
	}

	text, err := stringAt("json.decode", tok, args, 0)

	if err != nil {
		return err
	}

	var data interface{}

	if failure := json.Unmarshal([]byte(text), &data); failure != nil {
		return object.NewError(fault.Value, tok, "`json.decode()` cannot read this as JSON: %s", failure)
	}

	switch data := data.(type) {
	case []interface{}:
		elements := make([]object.Object, len(data))

		for index, value := range data {
			elements[index] = object.AnyValueToObject(value)
		}

		return &object.List{Elements: elements}
	case map[string]interface{}:
		// Go's json.Unmarshal decodes an object into a bare map, which does
		// not preserve the source text's key order, so neither can this -
		// the result is still frozen and consistent for every read after
		// this (§13.5), just not meaningfully "the order they appeared in
		// the JSON text".
		mapObject := object.NewOrderedMap()

		for key, value := range data {
			pairKey := &object.String{Value: key}

			mapObject.SetPair(pairKey.MapKey(), object.MapPair{Key: pairKey, Value: object.AnyValueToObject(value)})
		}

		return mapObject
	}

	// Valid JSON that is a bare number, string, boolean, or null. Decoding one
	// is not a failure of the text, so the message says what it found rather
	// than claiming the JSON is malformed.
	return object.NewError(fault.Value, tok, "`json.decode()` expects a JSON object or array at the top level").
		WithHelp("wrap the value in `[` and `]` to decode it as a list")
}

// jsonEncode renders a list or a map as JSON text.
func jsonEncode(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("json.encode", tok, args, 1); err != nil {
		return err
	}

	switch value := args[0].(type) {
	case *object.List:
		elements := make([]interface{}, len(value.Elements))

		for index, element := range value.Elements {
			elements[index] = object.ObjectToAnyValue(element)
		}

		return encodeJson(tok, elements)
	case *object.Map:
		// Reads value.Pairs directly rather than its insertion-ordered
		// OrderedPairs() (§13.5) - encodeJson's own json.Marshal call
		// always sorts a map's keys alphabetically regardless of the order
		// they were added to pairs here, so there is no order for this
		// loop to preserve or lose either way.
		pairs := make(map[string]interface{}, len(value.Pairs))

		for _, pair := range value.Pairs {
			key, ok := jsonKey(pair.Key)

			if !ok {
				return object.NewError(fault.Type, tok, "`json.encode()` cannot use %s as an object key", object.TypeName(pair.Key)).
					WithHelp("a JSON object key has to come from a string, a number, or a boolean")
			}

			pairs[key] = object.ObjectToAnyValue(pair.Value)
		}

		return encodeJson(tok, pairs)
	}

	return object.NewError(fault.Argument, tok, "`json.encode()` expects argument 1 to be a list or a map, got %s", object.TypeName(args[0]))
}

// encodeJson runs the Go encoder and turns whatever it objects to into a Ghost
// error, so a value the encoder cannot represent is reported rather than
// returned as an empty string.
func encodeJson(tok token.Token, value interface{}) object.Object {
	data, failure := json.Marshal(value)

	if failure != nil {
		return object.NewError(fault.Value, tok, "`json.encode()` cannot encode this value: %s", failure)
	}

	return &object.String{Value: string(data)}
}

// jsonKey renders a map key as the string a JSON object needs. JSON keys are
// always strings, so a number or a boolean key is written out as one.
func jsonKey(key object.Object) (string, bool) {
	switch key := key.(type) {
	case *object.String:
		return key.Value, true
	case *object.Number:
		return key.String(), true
	case *object.Boolean:
		return fmt.Sprintf("%t", key.Value), true
	}

	return "", false
}
