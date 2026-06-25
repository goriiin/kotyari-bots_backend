package logger

import (
	"errors"
	"testing"
)

// TestJoinErrs ensures the variadic logging helpers no longer silently drop
// errors when more than one is passed: joinErrs must surface every non-nil
// error and return nil only when there are none.
func TestJoinErrs(t *testing.T) {
	err1 := errors.New("first")
	err2 := errors.New("second")

	t.Run("nil when empty", func(t *testing.T) {
		if joinErrs() != nil {
			t.Fatal("expected nil for no errors")
		}
	})

	t.Run("nil when all nil", func(t *testing.T) {
		if joinErrs(nil, nil) != nil {
			t.Fatal("expected nil when all errors are nil")
		}
	})

	t.Run("single error passes through", func(t *testing.T) {
		got := joinErrs(err1)
		if !errors.Is(got, err1) {
			t.Fatalf("expected joined error to wrap err1, got %v", got)
		}
	})

	t.Run("multiple errors are all retained", func(t *testing.T) {
		got := joinErrs(err1, err2)
		if got == nil {
			t.Fatal("expected non-nil joined error")
		}
		if !errors.Is(got, err1) || !errors.Is(got, err2) {
			t.Fatalf("expected joined error to wrap both err1 and err2, got %v", got)
		}
	})
}
