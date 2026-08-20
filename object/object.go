package object

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

var evaluator func(node ast.Node, scope *Scope) Object

// Type identifies the runtime type of an object. It is an integer rather than a
// string so that the type comparisons the evaluator performs on every operation
// are single-word compares instead of string compares, and so that MapKey hashes
// over a machine word. String() keeps the human-readable name for error messages.
type Type int

const (
	BOOLEAN Type = iota
	BREAK
	CLASS
	CONTINUE
	ERROR
	FUNCTION
	INSTANCE
	LIBRARY_FUNCTION
	LIBRARY_MODULE
	LIBRARY_PROPERTY
	LIST
	MAP
	NULL
	NUMBER
	RETURN
	SCOPE
	STRING
	SUPER
	TRAIT
)

var typeNames = [...]string{
	BOOLEAN:          "BOOLEAN",
	BREAK:            "BREAK",
	CLASS:            "CLASS",
	CONTINUE:         "CONTINUE",
	ERROR:            "ERROR",
	FUNCTION:         "FUNCTION",
	INSTANCE:         "INSTANCE",
	LIBRARY_FUNCTION: "LIBRARY_FUNCTION",
	LIBRARY_MODULE:   "LIBRARY_MODULE",
	LIBRARY_PROPERTY: "LIBRARY_PROPERTY",
	LIST:             "LIST",
	MAP:              "MAP",
	NULL:             "NULL",
	NUMBER:           "NUMBER",
	RETURN:           "RETURN",
	SCOPE:            "SCOPE",
	STRING:           "STRING",
	SUPER:            "SUPER",
	TRAIT:            "TRAIT",
}

// String returns the name of the type, as used in runtime error messages.
func (t Type) String() string {
	if int(t) < 0 || int(t) >= len(typeNames) {
		return "UNKNOWN"
	}

	return typeNames[t]
}

// Object is the interface for all object values.
type Object interface {
	HasMethods
	Type() Type
	String() string
}

type MapKey struct {
	Type  Type
	Value uint64
}

type Mappable interface {
	MapKey() MapKey
}

type HasMethods interface {
	Method(method string, args []Object) (Object, bool)
}

type GoFunction func(scope *Scope, tok token.Token, args ...Object) Object
type GoProperty func(scope *Scope, tok token.Token) Object
type ObjectMethod func(value interface{}, args ...Object) (Object, bool)

func RegisterEvaluator(e func(node ast.Node, scope *Scope) Object) {
	evaluator = e
}

func AnyValueToObject(val any) Object {
	switch v := val.(type) {
	case bool:
		if v {
			return &Boolean{Value: true}
		}

		return &Boolean{Value: false}
	case string:
		return &String{Value: v}
	case int:
		return NewInt(int64(v))
	case int64:
		return NewInt(v)
	case float64:
		return NewFloat(v)
	case nil:
		return &Null{}
	case []any:
		elements := make([]Object, len(v))

		for index, item := range v {
			elements[index] = AnyValueToObject(item)
		}

		return &List{Elements: elements}
	case map[string]any:
		pairs := make(map[MapKey]MapPair)

		for key, val := range v {
			pairKey := &String{Value: key}
			var pairValue Object
			hashed := pairKey.MapKey()

			pairValue = AnyValueToObject(val)

			pairs[hashed] = MapPair{Key: pairKey, Value: pairValue}
		}

		return &Map{Pairs: pairs}
	}

	return nil
}

func ObjectToAnyValue(val Object) any {
	switch v := val.(type) {
	case *Boolean:
		return bool(v.Value)
	case *String:
		return string(v.Value)
	case *Number:
		if v.IsFloat() {
			return v.Float64()
		}
		return int(v.Int64())
	case *Null:
		return nil
	case *List:
		var collection []any

		for _, val := range v.Elements {
			collection = append(collection, ObjectToAnyValue(val))
		}

		return collection
	case *Map:
		collection := make(map[string]any)

		for _, pair := range v.Pairs {
			collection[pair.Key.(*String).Value] = ObjectToAnyValue(pair.Value)
		}

		return collection
	}

	return nil
}
