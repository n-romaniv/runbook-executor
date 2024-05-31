package statemachine

import (
	"encoding/json"
	"fmt"

	"github.com/n-romaniv/runbook-executor/internal/input"
	"github.com/n-romaniv/runbook-executor/internal/registry"
	"github.com/n-romaniv/runbook-executor/plugin"
)

type Definition struct {
	InitialState string                     `json:"initialState"`
	States       map[string]StateDefinition `json:"states"`
	Inputs       input.StateInput           `json:"inputs,omitempty"`
}

type StateDefinition struct {
	Action      string `json:"action,omitempty"`
	Transitions map[string]TransitionDefinition
	Inputs      input.StateInput `json:"inputs,omitempty"`
}

type TransitionDefinition struct {
	Next   string           `json:"next"`
	Inputs input.StateInput `json:"inputs,omitempty"`
}

func Parse(data []byte) (*StateMachine, error) {
	var def Definition
	err := json.Unmarshal(data, &def)
	if err != nil {
		return nil, err
	}

	// Put states in a map, then link them up
	states := make(map[string]*State, len(def.States))
	initializers := make([]plugin.MatcherWithInit, 0)

	for name, state := range def.States {
		s, err := buildState(name, state, def.Inputs)
		if err != nil {
			return nil, err
		}
		states[name] = s
	}

	for name, state := range def.States {
		for matcher, t := range state.Transitions {
			m, err := registry.LookupMatcher(matcher)
			if err != nil {
				return nil, err
			}

			if initializer, ok := m.(plugin.MatcherWithInit); ok {
				initializers = append(initializers, initializer)
			}

			next, ok := states[t.Next]
			cur := states[name]
			if !ok {
				return nil, fmt.Errorf("state %q from transition %q does not exist", t.Next, matcher)
			}
			cur.transitions = append(cur.transitions, &Transition{m, t.Inputs.Merged(cur.input), next})
		}
	}

	return &StateMachine{def.Inputs, states[def.InitialState], states, initializers}, nil
}

func buildState(name string, def StateDefinition, globalInputs input.StateInput) (*State, error) {
	var action plugin.Action
	var err error

	if def.Action != "" {
		action, err = registry.LookupAction(def.Action)
		if err != nil {
			return nil, err
		}
	}

	transitions := make([]*Transition, 0, len(def.Transitions))
	stateInputs := def.Inputs.Merged(globalInputs)

	return &State{name, stateInputs, action, transitions}, nil
}
