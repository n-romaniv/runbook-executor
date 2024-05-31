package core

import (
	"context"
	"errors"
	"time"

	"github.com/n-romaniv/runbook-executor/internal/input"
	"github.com/n-romaniv/runbook-executor/internal/registry"
	"github.com/n-romaniv/runbook-executor/plugin"
)

const defaultWait = 30

func Not(ctx context.Context, in plugin.InputProvider) (<-chan struct{}, error) {
	i, ok := in.(input.StateInput)
	if !ok {
		return nil, errors.New("can only use not combinator with state input")
	}

	m := in.GetString("matcher")
	s := in.GetInt("for")
	subCfg := in.Get("inputs")
	if m == "" {
		return nil, errors.New("matcher input can not be empty")
	}
	if s == 0 {
		s = defaultWait
	}

	matcher, err := registry.LookupMatcher(m)
	if err != nil {
		return nil, err
	}

	var match <-chan struct{}
	subInput, ok := subCfg.(input.StateInput)
	cancelCtx, cancel := context.WithCancel(ctx)
	if ok {
		// TODO: It should probably only have stuff from the global inputs, but inputs of Times combinator.
		match, err = matcher.Wait(cancelCtx, subInput.Merged(i))
	} else {
		match, err = matcher.Wait(cancelCtx, in)
	}

	if err != nil {
		cancel()
		return nil, err
	}

	output := make(chan struct{})

	go func() {
		defer cancel()
		select {
		case <-match:
			return
		case <-time.After(time.Duration(s) * time.Second):
			output <- struct{}{}
		}
	}()

	return output, nil
}

func Times(ctx context.Context, in plugin.InputProvider) (<-chan struct{}, error) {
	i, ok := in.(input.StateInput)
	if !ok {
		return nil, errors.New("can only use times combinator with state input")
	}

	m := in.GetString("matcher")
	n := in.GetInt("n")
	subCfg := in.Get("inputs")
	if m == "" {
		return nil, errors.New("matcher input can not be empty")
	}
	if n == 0 {
		return nil, errors.New("input n is required")
	}

	matcher, err := registry.LookupMatcher(m)
	if err != nil {
		return nil, err
	}

	var match <-chan struct{}
	subInput, ok := subCfg.(input.StateInput)
	cancelCtx, cancel := context.WithCancel(ctx)
	// TODO: It should probably only have stuff from the global inputs, but inputs of Times combinator.
	if ok {
		match, err = matcher.Wait(cancelCtx, subInput.Merged(i))
	} else {
		match, err = matcher.Wait(cancelCtx, in)
	}

	if err != nil {
		cancel()
		return nil, err
	}

	output := make(chan struct{})
	counter := 0

	go func() {
		defer cancel()
		for range match {
			counter += 1
			if counter >= n {
				output <- struct{}{}
				return
			}
		}
	}()

	return output, nil
}

// TODO: Add things like And, Or etc.
