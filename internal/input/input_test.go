package input

import (
	"maps"
	"slices"
	"testing"
)

func TestGetStringSlice(t *testing.T) {
	input := StateInput{
		"exists":    []string{"apple", "banana"},
		"wrongType": "not a slice",
	}

	tests := []struct {
		name     string
		key      string
		expected []string
	}{
		{"Key exists", "exists", []string{"apple", "banana"}},
		{"Key type mismatch", "wrongType", []string{}},
		{"Key does not exist", "doesNotExist", []string{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := input.GetStringSlice(tc.key); !slices.Equal(got, tc.expected) {
				t.Errorf("GetStringSlice(%q) = %v, want %v", tc.key, got, tc.expected)
			}
		})
	}
}

func TestGetString(t *testing.T) {
	input := StateInput{
		"valid":   "hello",
		"invalid": 123,
	}

	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{"Valid key", "valid", "hello"},
		{"Invalid type", "invalid", ""},
		{"Non-existent key", "notHere", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := input.GetString(tc.key); got != tc.expected {
				t.Errorf("GetString(%q) = %v, want %v", tc.key, got, tc.expected)
			}
		})
	}
}

func TestMerged(t *testing.T) {
	input1 := StateInput{
		"name":  "example",
		"count": 10,
	}
	input2 := StateInput{
		"count":  20,
		"newKey": "newValue",
	}
	expected := StateInput{
		"name": "example",
		// Should not be overwritten
		"count":  10,
		"newKey": "newValue",
	}

	result := input1.Merged(input2)
	if !maps.Equal(result, expected) {
		t.Errorf("Merged result = %v, want %v", result, expected)
	}
}
