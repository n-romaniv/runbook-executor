package statemachine

import (
	"context"
	"errors"
	"reflect"

	"github.com/n-romaniv/runbook-executor/internal/input"
	"github.com/n-romaniv/runbook-executor/plugin"
)

type StateMachine struct {
	input        input.StateInput
	initial      *State
	states       map[string]*State
	initializers []plugin.MatcherWithInit
}

type State struct {
	name        string
	input       input.StateInput
	action      plugin.Action
	transitions []*Transition
}

type Transition struct {
	matcher plugin.Matcher
	input   input.StateInput
	next    *State
}

func (s *State) Name() string {
	return s.name
}

func (m *StateMachine) Run(ctx context.Context) <-chan error {
	doneChan := make(chan error)

	for _, i := range m.initializers {
		i.Init(ctx, m.input)
	}

	m.initial.Run(ctx, doneChan)

	return doneChan
}

func (m *StateMachine) Validate() error {
	visited := map[*State]struct{}{m.initial: {}}

	traverse(m.initial, visited)

	if len(visited) != len(m.states) {
		return errors.New("some states are unreachable")
	}

	return nil
}

func traverse(s *State, visited map[*State]struct{}) {
	visited[s] = struct{}{}

	for _, t := range s.transitions {
		if _, ok := visited[t.next]; !ok {
			traverse(t.next, visited)
		}
	}
}

func (s *State) Run(ctx context.Context, doneChan chan<- error) {
	err := s.action.Execute(ctx, s.input)

	if err != nil {
		doneChan <- err
		return
	}

	cases := make([]reflect.SelectCase, 0, len(s.transitions))
	cancelCtx, cancel := context.WithCancel(ctx)

	if len(s.transitions) == 0 {
		close(doneChan)
		cancel()
		return
	}

	for _, t := range s.transitions {
		ch, err := t.matcher.Wait(cancelCtx, s.input)

		if err != nil {
			cancel()
			doneChan <- err
			return
		}

		cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)})
	}

	go func() {
		defer cancel()

		chosen, _, _ := reflect.Select(cases)
		s.transitions[chosen].next.Run(ctx, doneChan)
	}()
}
