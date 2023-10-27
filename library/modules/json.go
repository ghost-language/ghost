package modules

import (
	"encoding/json"
	"fmt"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

var JsonMethods = map[string]*object.LibraryFunction{}
var JsonProperties = map[string]*object.LibraryProperty{}

func init() {
	RegisterMethod(JsonMethods, "decode", jsonDecode)
	RegisterMethod(JsonMethods, "encode", jsonEncode)
}

// jsonDecode decodes the JSON-encoded data and returns a new list or map object.
func jsonDecode(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("wrong number of arguments. got=%d, want=1", len(args))
	}

	str, ok := args[0].(*object.String)

	if !ok {
		return object.NewError("argument to `decode` must be STRING, got %s", args[0].Type())
	}

	var data interface{}

	err := json.Unmarshal([]byte(str.Value), &data)

	if err != nil {
		return object.NewError("failed to decode JSON: %s", err.Error())
	}

	switch v := data.(type) {
	case []interface{}:
		var elements []object.Object

		for _, val := range v {
			elements = append(elements, object.AnyValueToObject(val))
		}

		return &object.List{Elements: elements}
	case map[string]interface{}:
		pairs := make(map[object.MapKey]object.MapPair)

		for key, val := range v {
			pairKey := &object.String{Value: key}
			pairValue := object.AnyValueToObject(val)

			pairs[pairKey.MapKey()] = object.MapPair{Key: pairKey, Value: pairValue}
		}

		return &object.Map{Pairs: pairs}
	}

	return object.NewError("failed to decode JSON: %s", err.Error())
}

// jsonEncode returns the JSON encoding of either a list or map object.
func jsonEncode(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.List:
		var elements []interface{}

		for _, val := range arg.Elements {
			elements = append(elements, object.ObjectToAnyValue(val))
		}

		data, err := json.Marshal(elements)

		if err != nil {
			return object.NewError("failed to encode JSON: %s", err.Error())
		}

		return &object.String{Value: string(data)}
	case *object.Map:
		pairs := make(map[string]interface{})

		for _, pair := range arg.Pairs {
			// map keys can be numbers, strings, or booleans
			switch pair.Key.(type) {
			case *object.String:
				pairs[pair.Key.(*object.String).Value] = object.ObjectToAnyValue(pair.Value)
			case *object.Number:
				pairs[fmt.Sprintf("%d", object.ObjectToAnyValue(pair.Key.(*object.Number)))] = object.ObjectToAnyValue(pair.Value)
			case *object.Boolean:
				pairs[fmt.Sprintf("%t", pair.Key.(*object.Boolean).Value)] = object.ObjectToAnyValue(pair.Value)
			}
		}

		data, err := json.Marshal(pairs)

		if err != nil {
			return object.NewError("failed to encode JSON: %s", err.Error())
		}

		return &object.String{Value: string(data)}
	}

	return object.NewError("argument to `encode` must be LIST or MAP, got %s", args[0].Type())
}
