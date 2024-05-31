package core

import (
	"context"
	"testing"
	"time"

	"github.com/n-romaniv/runbook-executor/internal/input"
	"github.com/n-romaniv/runbook-executor/plugin"
)

// Example check functions
func alwaysTrue(_ context.Context, _ plugin.InputProvider) bool {
	return true
}

func alwaysFalse(_ context.Context, _ plugin.InputProvider) bool {
	return false
}
func TestSimpleMatcher_Wait_SignalEmitted(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Small interval for test speed-up
	matcher := NewSimpleMatcher(alwaysTrue, 10*time.Millisecond)

	ch, err := matcher.Wait(ctx, input.StateInput{})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	select {
	case <-ch:
		// Signal received as expected
	case <-time.After(50 * time.Millisecond):
		t.Fatal("Expected a signal but none was received")
	}
}

func TestSimpleMatcher_Wait_NoSignalEmitted(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	matcher := NewSimpleMatcher(alwaysFalse, 10*time.Millisecond)

	ch, err := matcher.Wait(ctx, input.StateInput{})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	select {
	case <-ch:
		t.Fatal("Did not expect a signal but one was received")
	case <-time.After(50 * time.Millisecond):
		// No signal received as expected
	}
}

func TestSimpleMatcher_Wait_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	matcher := NewSimpleMatcher(alwaysTrue, 10*time.Millisecond)

	ch, err := matcher.Wait(ctx, input.StateInput{})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	select {
	case <-ch:
		t.Fatal("Did not expect a signal after context cancellation")
	case <-time.After(50 * time.Millisecond):
		// Test passed, no signal received, context cancellation observed
	}
}
