package modules

import (
	"testing"

	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

func TestOsSleep(t *testing.T) {
	result := osSleep(nil, token.Token{}, object.NewInt(1))

	if object.IsError(result) {
		t.Fatalf("unexpected error: %s", result.String())
	}
}

func TestOsSleepRejectsNegativeDuration(t *testing.T) {
	result := osSleep(nil, token.Token{}, object.NewInt(-1))

	if !object.IsError(result) {
		t.Fatalf("expected an error for a negative duration, got=%v", result)
	}
}

func TestOsSleepArity(t *testing.T) {
	result := osSleep(nil, token.Token{})

	if !object.IsError(result) {
		t.Fatalf("expected an error for missing arguments, got=%v", result)
	}
}
