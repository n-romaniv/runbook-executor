package statemachine

import (
	"testing"
)

func TestParseValidStateMachine(t *testing.T) {
	jsonInput := `{
        "initialState": "start",
        "states": {
            "start": {
                "transitions": {
                    "always": {
                        "next": "end"
                    }
                }
            },
            "end": {}
        }
    }`

	expectedInitialState := "start"
	sm, err := Parse([]byte(jsonInput))
	if err != nil {
		t.Fatalf("Parse failed with error: %v", err)
	}
	if sm.initial.name != expectedInitialState {
		t.Errorf("Expected initial state %s, got %s", expectedInitialState, sm.initial.name)
	}
}
