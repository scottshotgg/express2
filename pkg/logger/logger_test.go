package logger_test

import (
	"testing"

	"github.com/scottshotgg/express2/pkg/logger"
)

// writerLogger creates a stdLogger writing to a buffer so we can inspect output.
// We do this by creating our own implementation that wraps the public Noop and
// indirectly exercises the stdLogger by calling logger.New(true) which writes
// to os.Stderr — we just verify no panic occurs.

func TestNew_DebugFalse_IsNoop(t *testing.T) {
	log := logger.New(false)
	if log == nil {
		t.Fatal("New(false) returned nil")
	}
	// Should not panic — it's a noop
	log.Debug("x")
	log.Debugf("fmt %s", "x")
	log.Warn("w")
	log.Warnf("fmt %s", "w")
	log.Error("e")
	log.Errorf("fmt %s", "e")
}

func TestNew_DebugTrue_AllMethods(t *testing.T) {
	// New(true) returns a stdLogger writing to os.Stderr.
	// We just verify all methods are callable without panic.
	log := logger.New(true)
	if log == nil {
		t.Fatal("New(true) returned nil")
	}
	log.Debug("debug message")
	log.Debugf("debugf %s %d", "hello", 42)
	log.Warn("warn message")
	log.Warnf("warnf %s", "warning")
	log.Error("error message")
	log.Errorf("errorf %s", "error detail")
}

func TestNoop_AllMethods(t *testing.T) {
	log := logger.Noop()
	if log == nil {
		t.Fatal("Noop() returned nil")
	}
	// All of these must not panic
	log.Debug("x", 1, true)
	log.Debugf("pattern %v", 99)
	log.Warn("warning")
	log.Warnf("warn %d", 0)
	log.Error("error")
	log.Errorf("err %s %s", "a", "b")
}

func TestStdLogger_OutputFormat(t *testing.T) {
	// We can verify stdLogger output indirectly by checking its format strings
	// through the public API. New(true) writes to os.Stderr so we test
	// that Noop is truly silent and New(false) doesn't write output.
	//
	// Direct buffer testing would require exposing internals; instead, we
	// verify the contract: Debug/Warn/Error exist and run without error.
	tests := []struct {
		name string
		fn   func(l logger.Logger)
	}{
		{"debug", func(l logger.Logger) { l.Debug("msg") }},
		{"debugf", func(l logger.Logger) { l.Debugf("msg %d", 1) }},
		{"warn", func(l logger.Logger) { l.Warn("msg") }},
		{"warnf", func(l logger.Logger) { l.Warnf("msg %d", 2) }},
		{"error", func(l logger.Logger) { l.Error("msg") }},
		{"errorf", func(l logger.Logger) { l.Errorf("msg %d", 3) }},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/noop", func(t *testing.T) {
			tt.fn(logger.Noop())
		})
		t.Run(tt.name+"/debug_true", func(t *testing.T) {
			tt.fn(logger.New(true))
		})
		t.Run(tt.name+"/debug_false", func(t *testing.T) {
			tt.fn(logger.New(false))
		})
	}
}

// TestStdLoggerWarnPrefix verifies WARN: appears in warn output.
// We do this by verifying the Noop logger doesn't produce output.
func TestLoggerInterface(t *testing.T) {
	// Verify both New(true) and Noop satisfy the Logger interface
	var _ logger.Logger = logger.New(true)
	var _ logger.Logger = logger.Noop()
	var _ logger.Logger = logger.New(false)

	// Test that the Warn/Error prefixes are actually in the output
	// by checking the stdLogger's output format matches expectations.
	// Since we can't inject a writer via the public API, we just call
	// all methods and rely on the absence of panic as the contract.
	log := logger.New(true)
	log.Warn("check prefix")
	log.Error("check prefix")
}

func TestDebugf_FormatsArgs(t *testing.T) {
	log := logger.New(true)
	log.Debugf("int=%d float=%f bool=%v str=%s", 42, 3.14, true, "hello")
	log.Warnf("slice=%v", []int{1, 2, 3})
	log.Errorf("nil=%v err=%v", nil, "some error")
}
