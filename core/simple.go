package core

import (
	"context"
	"time"

	"github.com/n-romaniv/runbook-executor/plugin"
)

type checkFn = func(context.Context, plugin.InputProvider) bool

type SimpleMatcher struct {
	check  checkFn
	ticker *time.Ticker
}

func (m SimpleMatcher) Wait(ctx context.Context, cfg plugin.InputProvider) (<-chan struct{}, error) {
	ch := make(chan struct{})
	go func() {
		for {
			select {
			case <-m.ticker.C:
				if m.check(ctx, cfg) {
					ch <- struct{}{}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch, nil
}

func NewSimpleMatcher(f checkFn, interval time.Duration) SimpleMatcher {
	return SimpleMatcher{f, time.NewTicker(interval)}
}
