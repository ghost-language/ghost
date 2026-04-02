package modules

import (
	"time"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

var TimeMethods = map[string]*object.LibraryFunction{}
var TimeProperties = map[string]*object.LibraryProperty{}

func init() {
	RegisterMethod(TimeMethods, "sleep", timeSleep)
	RegisterMethod(TimeMethods, "now", timeNow)

	RegisterProperty(TimeProperties, "nanosecond", timeNanosecond)
	RegisterProperty(TimeProperties, "microsecond", timeMicrosecond)
	RegisterProperty(TimeProperties, "millisecond", timeMillisecond)
	RegisterProperty(TimeProperties, "second", timeSecond)
	RegisterProperty(TimeProperties, "minute", timeMinute)
	RegisterProperty(TimeProperties, "hour", timeHour)
	RegisterProperty(TimeProperties, "day", timeDay)
	RegisterProperty(TimeProperties, "week", timeWeek)
	RegisterProperty(TimeProperties, "month", timeMonth)
	RegisterProperty(TimeProperties, "year", timeYear)
}

func timeSleep(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 1 {
		return nil
	}

	if args[0].Type() != object.NUMBER {
		return nil
	}

	ms := args[0].(*object.Number)
	time.Sleep(time.Duration(ms.Int64()) * time.Millisecond)

	return nil
}

func timeNow(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if len(args) != 0 {
		return nil
	}

	return object.NewInt(time.Now().Unix())
}

// properties

func timeNanosecond(scope *object.Scope, tok token.Token) object.Object {
	return object.NewFloat(0.00001)
}

func timeMicrosecond(scope *object.Scope, tok token.Token) object.Object {
	return object.NewFloat(0.0001)
}

func timeMillisecond(scope *object.Scope, tok token.Token) object.Object {
	return object.NewFloat(0.001)
}

func timeSecond(scope *object.Scope, tok token.Token) object.Object {
	return object.NewInt(1)
}

func timeMinute(scope *object.Scope, tok token.Token) object.Object {
	return object.NewInt(60)
}

func timeHour(scope *object.Scope, tok token.Token) object.Object {
	return object.NewInt(3600)
}

func timeDay(scope *object.Scope, tok token.Token) object.Object {
	return object.NewInt(86400)
}

func timeWeek(scope *object.Scope, tok token.Token) object.Object {
	return object.NewInt(604800)
}

func timeMonth(scope *object.Scope, tok token.Token) object.Object {
	return object.NewInt(2592000)
}

func timeYear(scope *object.Scope, tok token.Token) object.Object {
	return object.NewInt(31536000)
}
