package statemachine

import (
	"context"
	"testing"
	"time"

	"github.com/n-romaniv/runbook-executor/core"
	"github.com/n-romaniv/runbook-executor/internal/input"
	"github.com/n-romaniv/runbook-executor/plugin"
)

func TestValidate(t *testing.T) {
	// Create states
	stateA := &State{name: "A", input: input.StateInput{}, transitions: []*Transition{}}
	stateB := &State{name: "B", input: input.StateInput{}, transitions: []*Transition{}}
	stateC := &State{name: "C", input: input.StateInput{}, transitions: []*Transition{}}

	// Link states
	stateA.transitions = append(stateA.transitions, &Transition{next: stateB})

	sm := &StateMachine{
		initial: stateA,
		states: map[string]*State{
			"A": stateA,
			"B": stateB,
			"C": stateC, // State C is unreachable
		},
	}

	// Test validation
	err := sm.Validate()
	if err == nil || err.Error() != "some states are unreachable" {
		t.Errorf("Validate should have identified unreachable states, got: %v", err)
	}
}

func TestRun_NoTransitions(t *testing.T) {
	ctx := context.Background()

	initial := &State{name: "initial", input: input.StateInput{}, transitions: []*Transition{}}
	sm := &StateMachine{
		initial: initial,
		states:  map[string]*State{"initial": initial},
	}

	doneChan := sm.Run(ctx)
	select {
	case err := <-doneChan:
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Expected to finish without delay")
	}
}

func TestRun_SuccessfulTransition(t *testing.T) {
	ctx := context.Background()
	transitioned := false

	nextAction := core.ActionFn(func(ctx context.Context, inputs plugin.InputProvider) error {
		transitioned = true
		return nil
	})

	next := &State{name: "next", input: input.StateInput{}, action: nextAction, transitions: []*Transition{}}
	initial := &State{
		name:  "initial",
		input: input.StateInput{},
		transitions: []*Transition{
			{matcher: core.Always, input: input.StateInput{}, next: next},
		},
	}

	sm := &StateMachine{
		initial: initial,
		states:  map[string]*State{"initial": initial, "next": next},
	}

	doneChan := sm.Run(ctx)
	err := <-doneChan
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !transitioned {
		t.Errorf("Expected to transition to 'next' state")
	}
}
