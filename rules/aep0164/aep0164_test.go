package aep0164

import (
	"testing"

	"github.com/aep-dev/api-linter/lint"
)

func TestAddRules(t *testing.T) {
	if err := AddRules(lint.NewRuleRegistry()); err != nil {
		t.Errorf("AddRules got an error: %v", err)
	}
}
