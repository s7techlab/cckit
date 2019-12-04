package matchers

import (
	"fmt"
	"strings"

	"github.com/onsi/gomega/format"
)

type ErrorIslMatcher struct {
	Expected interface{}
}

func (matcher *ErrorIslMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil && matcher.Expected == nil {
		return false, fmt.Errorf("Refusing to compare <nil> to <nil>.\nBe explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized.")
	}

	// 1.12
	return strings.Contains(
		fmt.Sprintf(`%s`, actual),
		fmt.Sprintf(`%s`, matcher.Expected)), nil
}

func (matcher *ErrorIslMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to match error", matcher.Expected)
}

func (matcher *ErrorIslMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to match error", matcher.Expected)
}
