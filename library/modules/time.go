package modules

import (
	"time"

	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

var TimeMethods = map[string]*object.LibraryFunction{}
var TimeProperties = map[string]*object.LibraryProperty{}

func init() {
	RegisterMethod(TimeMethods, "sleep", timeSleep)
	RegisterMethod(TimeMethods, "now", timeNow)
	RegisterMethod(TimeMethods, "nowNano", timeNowNano)

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
	if err := arity("time.sleep", tok, args, 1); err != nil {
		return err
	}

	milliseconds, err := integerAt("time.sleep", tok, args, 0)

	if err != nil {
		return err
	}

	if milliseconds < 0 {
		return object.NewError(fault.Value, tok, "`time.sleep()` expects a duration of zero or greater, got %d", milliseconds)
	}

	time.Sleep(time.Duration(milliseconds) * time.Millisecond)

	return nil
}

func timeNow(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("time.now", tok, args, 0); err != nil {
		return err
	}

	return object.NewInt(time.Now().Unix())
}

// timeNowNano is time.now() at nanosecond precision, for measuring how long
// something took rather than what time it is. It replaces os.clock(), which
// answered the same question from a module about the operating system rather
// than about time.
func timeNowNano(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("time.nowNano", tok, args, 0); err != nil {
		return err
	}

	return object.NewInt(time.Now().UnixNano())
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
