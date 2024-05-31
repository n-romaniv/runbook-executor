package plugin

import (
	"context"
)

func (b *PluginDefinition) WithInit(f func(context.Context) error) *PluginDefinition {
	b.init = f
	return b
}

func String(name string, optional bool) inputRequirement {
	return inputRequirement{name, stringInput, optional}
}

func StringSlice(name string, optional bool) inputRequirement {
	return inputRequirement{name, stringSliceInput, optional}
}

func Int(name string, optional bool) inputRequirement {
	return inputRequirement{name, intInput, optional}
}

func Float(name string, optional bool) inputRequirement {
	return inputRequirement{name, floatInput, optional}
}

func Bool(name string, optional bool) inputRequirement {
	return inputRequirement{name, boolInput, optional}
}

func (b *PluginDefinition) Action(name string, a Action, inputs ...inputRequirement) *PluginDefinition {
	def := ActionDefinition{name, a, inputs}
	b.actions = append(b.actions, def)
	return b
}

func (b *PluginDefinition) Matcher(name string, m Matcher, inputs ...inputRequirement) *PluginDefinition {
	def := MatcherDefinition{name, m, inputs}
	b.matchers = append(b.matchers, def)
	return b
}

func New(name string) PluginDefinition {
	return PluginDefinition{name: name}
}

func (d PluginDefinition) Name() string {
	return d.name
}

func (d PluginDefinition) InitFunc() func(context.Context) error {
	return d.init
}

func (d PluginDefinition) Actions() []ActionDefinition {
	return d.actions
}

func (d PluginDefinition) Matchers() []MatcherDefinition {
	return d.matchers
}

func (a ActionDefinition) Action() Action {
	return a.action
}

type actionWithInputValidation struct {
	a      Action
	inputs []inputRequirement
}

type matcherWithInputValidation struct {
	m      Matcher
	inputs []inputRequirement
}

func ActionWithInputValidation(a Action, inputs []inputRequirement) ActionWithValidation {
	return actionWithInputValidation{a, inputs}
}

func MatcherWithInputValidation(m Matcher, inputs []inputRequirement) MatcherWithValidation {
	return matcherWithInputValidation{m, inputs}
}

func (a actionWithInputValidation) Execute(ctx context.Context, in InputProvider) error {
	return a.a.Execute(ctx, in)
}

func (a actionWithInputValidation) Validate(ctx context.Context, in InputProvider) error {
	// TODO: implement actual validation
	return nil
}

func (m matcherWithInputValidation) Wait(ctx context.Context, in InputProvider) (<-chan struct{}, error) {
	return m.m.Wait(ctx, in)
}

func (m matcherWithInputValidation) Validate(ctx context.Context, in InputProvider) error {
	// TODO: implement actual validation
	return nil
}
