package registry

import (
	"context"
	"fmt"
	"sync"

	"github.com/n-romaniv/runbook-executor/plugin"
)

var (
	actionsMu      sync.RWMutex
	matchersMu     sync.RWMutex
	initFuncsMu    sync.RWMutex
	actionPlugins  = make(map[string]plugin.Action)
	matcherPlugins = make(map[string]plugin.Matcher)
	initFuncs      = make(map[string]func(context.Context) error)
)

func RegisterAction(name string, a plugin.Action) {
	actionsMu.Lock()
	defer actionsMu.Unlock()
	if _, dup := actionPlugins[name]; dup {
		panic("Register called twice for action " + name)
	}
	actionPlugins[name] = a
}

func RegisterMatcher(name string, m plugin.Matcher) {
	matchersMu.Lock()
	defer matchersMu.Unlock()
	if _, dup := matcherPlugins[name]; dup {
		panic("Register called twice for matcher " + name)
	}
	matcherPlugins[name] = m
}

func RegisterInitAction(name string, f func(context.Context) error) {
	initFuncsMu.Lock()
	defer initFuncsMu.Unlock()
	if _, dup := matcherPlugins[name]; dup {
		panic("Register called twice for init function " + name)
	}
	initFuncs[name] = f
}

func LookupAction(name string) (plugin.Action, error) {
	actionsMu.RLock()
	defer actionsMu.RUnlock()
	a, ok := actionPlugins[name]

	if !ok {
		return nil, fmt.Errorf("action %q not found", name)
	}

	return a, nil
}

func LookupMatcher(name string) (plugin.Matcher, error) {
	matchersMu.RLock()
	defer matchersMu.RUnlock()
	a, ok := matcherPlugins[name]

	if !ok {
		return nil, fmt.Errorf("action %q not found", name)
	}

	return a, nil
}

func RegisterPlugin(plugin plugin.PluginDefinition) {
	// n := plugin.Name()

}
