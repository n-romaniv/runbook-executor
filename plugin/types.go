package plugin

import (
	"context"
)

type InputProvider interface {
	GetStringSlice(key string) []string
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetFloat64(key string) float64
	Get(key string) any
}

type Action interface {
	Execute(ctx context.Context, in InputProvider) error
}

type Matcher interface {
	Wait(ctx context.Context, in InputProvider) (<-chan struct{}, error)
}

type MatcherWithInit interface {
	Matcher
	Init(ctx context.Context, in InputProvider) error
}

type MatcherWithValidation interface {
	Matcher
	Validate(ctx context.Context, in InputProvider) error
}

type ActionWithValidation interface {
	Action
	Validate(ctx context.Context, in InputProvider) error
}

type inputType int

const (
	stringInput inputType = iota
	stringSliceInput
	intInput
	floatInput
	boolInput
)

type inputRequirement struct {
	name     string
	t        inputType
	optional bool
}

type ActionDefinition struct {
	name         string
	action       Action
	requirements []inputRequirement
}

type MatcherDefinition struct {
	name         string
	matcher      Matcher
	requirements []inputRequirement
}

type PluginDefinition struct {
	name     string
	init     func(context.Context) error
	actions  []ActionDefinition
	matchers []MatcherDefinition
}
