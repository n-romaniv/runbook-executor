package core

import (
	"context"
	"testing"
	"time"

	"github.com/n-romaniv/runbook-executor/internal/input"
	"github.com/n-romaniv/runbook-executor/internal/registry"
	"github.com/n-romaniv/runbook-executor/plugin"
)

var neverMatch = MatcherFn(func(ctx context.Context, cp plugin.InputProvider) (<-chan struct{}, error) {
	return make(<-chan struct{}), nil
})

func init() {
	registry.RegisterMatcher("never", neverMatch)
}

func TestNot_Triggered(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cfg := input.StateInput{
		"matcher": "never",
		"for":     2,
	}

	output, err := Not(ctx, cfg)
	if err != nil {
		t.Fatalf("Not should not have returned an error, got %v", err)
	}

	select {
	case <-output:
		// Passed: output was triggered as expected
	case <-time.After(2*time.Second + 10*time.Millisecond):
		t.Fatal("Expected output to be triggered after 2 seconds")
	}
}

func TestNot_NotTriggered(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cfg := input.StateInput{
		"matcher": "never",
		"for":     2,
	}

	output, err := Not(ctx, cfg)
	if err != nil {
		t.Fatalf("Not should not have returned an error, got %v", err)
	}

	select {
	case <-output:
		// Passed: output was triggered as expected
	case <-time.After(2*time.Second + 10*time.Millisecond):
		t.Fatal("Expected output to be triggered after 2 seconds")
	}
}

func TestTimes_Triggered(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	triggerTimes := 3
	countingMatcher := MatcherFn(func(ctx context.Context, cfg plugin.InputProvider) (<-chan struct{}, error) {
		ch := make(chan struct{}, triggerTimes)
		for range triggerTimes {
			ch <- struct{}{}
		}
		return ch, nil
	})

	registry.RegisterMatcher("counter-pass", countingMatcher)

	cfg := input.StateInput{
		"matcher": "counter-pass",
		"n":       triggerTimes,
	}

	output, err := Times(ctx, cfg)
	if err != nil {
		t.Fatalf("Times should not have returned an error, got %v", err)
	}

	select {
	case <-output:
		// Passed: output was triggered as expected after the specified number of triggers
	case <-time.After(1 * time.Second):
		t.Fatal("Expected output to be triggered after the specified number of matches")
	}
}

func TestTimes_NotTriggered(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	triggerTimes := 3
	countingMatcher := MatcherFn(func(ctx context.Context, cfg plugin.InputProvider) (<-chan struct{}, error) {
		ch := make(chan struct{}, triggerTimes)
		for range triggerTimes {
			ch <- struct{}{}
		}
		return ch, nil
	})

	registry.RegisterMatcher("counter-fail", countingMatcher)

	cfg := input.StateInput{
		"matcher": "counter-fail",
		"n":       triggerTimes + 1,
	}

	output, err := Times(ctx, cfg)
	if err != nil {
		t.Fatalf("Times should not have returned an error, got %v", err)
	}

	select {
	case <-output:
		t.Fatal("Did not expect output to be triggered after the specified number of matches")
	case <-time.After(1 * time.Second):
		// Should not be triggered within 1 second, all good
	}
}
