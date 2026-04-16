package builder_test

import (
	"strings"
	"testing"

	"github.com/scottshotgg/express2/builder"
)

func TestErrConstants(t *testing.T) {
	if builder.ErrNotImplemented == nil {
		t.Error("ErrNotImplemented is nil")
	}
	if builder.ErrMultDimArrInit == nil {
		t.Error("ErrMultDimArrInit is nil")
	}
	if builder.ErrOutOfTokens == nil {
		t.Error("ErrOutOfTokens is nil")
	}
}

func TestAppendTokenToError(t *testing.T) {
	b, err := getBuilderFromString("42")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	appErr := b.AppendTokenToError("test error")
	if appErr == nil {
		t.Fatal("AppendTokenToError returned nil")
	}
	// The error should contain our message
	if !strings.Contains(appErr.Error(), "test error") {
		t.Errorf("error %q does not contain 'test error'", appErr.Error())
	}
}

func TestAppendTokenToError_AtEnd(t *testing.T) {
	b, err := getBuilderFromString("42")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	// Advance to the last token
	for b.Index < len(b.Tokens)-1 {
		b.Index++
	}
	appErr := b.AppendTokenToError("end error")
	if appErr == nil {
		t.Fatal("AppendTokenToError returned nil at end")
	}
	if !strings.Contains(appErr.Error(), "no token to print") {
		t.Errorf("error %q should contain 'no token to print'", appErr.Error())
	}
}
