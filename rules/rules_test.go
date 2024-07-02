package rules

import (
	"errors"
	"testing"

	"github.com/aep-dev/api-linter/lint"
)

func TestAddAEPRules(t *testing.T) {
	wantError := errors.New("test")
	tests := []struct {
		name          string
		addRulesFuncs []addRulesFuncType
		err           error
	}{
		{
			name:          "EmptyRules_NoError",
			addRulesFuncs: nil,
			err:           nil,
		},
		{
			name: "AddingRules_NoError",
			addRulesFuncs: []addRulesFuncType{
				func(lint.RuleRegistry) error { return nil },
			},
			err: nil,
		},
		{
			name: "ReturnError",
			addRulesFuncs: []addRulesFuncType{
				func(lint.RuleRegistry) error { return wantError },
			},
			err: wantError,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := addAEPRules(lint.NewRuleRegistry(), test.addRulesFuncs)
			if err != test.err {
				t.Errorf("addAEPRules got error %v, but want %v", err, test.err)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	if err := Add(lint.NewRuleRegistry()); err != nil {
		t.Errorf("Add got an error: %v", err)
	}
}
