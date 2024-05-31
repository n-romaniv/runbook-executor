package core

import "github.com/n-romaniv/runbook-executor/internal/registry"

func init() {
	registry.RegisterMatcher("always", Always)
	registry.RegisterMatcher("not", MatcherFn(Not))
	registry.RegisterMatcher("times", MatcherFn(Times))
}
