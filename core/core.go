package core

import (
	"context"

	"github.com/n-romaniv/runbook-executor/plugin"
)

type ActionFn func(context.Context, plugin.InputProvider) error
type MatcherFn func(context.Context, plugin.InputProvider) (<-chan struct{}, error)

func (a ActionFn) Execute(ctx context.Context, cfg plugin.InputProvider) error {
	return a(ctx, cfg)
}

func (m MatcherFn) Wait(ctx context.Context, cfg plugin.InputProvider) (<-chan struct{}, error) {
	return m(ctx, cfg)
}

func always(context.Context, plugin.InputProvider) (<-chan struct{}, error) {
	ch := make(chan struct{})
	close(ch)
	return ch, nil
}

var Always = MatcherFn(always)
