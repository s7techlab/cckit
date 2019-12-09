package matchers

import (
	"errors"
	"fmt"

	"github.com/onsi/gomega/format"
)

type StringerEqualMatcher struct {
	Expected interface{}
}

type Stringer interface {
	String() string
}

func (matcher *StringerEqualMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil && matcher.Expected == nil {
		return false, fmt.Errorf("Refusing to compare <nil> to <nil>.\nBe explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized.")
	}

	actualStringer, okActual := actual.(Stringer)
	if !okActual {
		return false, errors.New("refusing to compare non-stringer actual value")
	}

	expectedStringer, okExpected := matcher.Expected.(Stringer)
	if !okExpected {
		return false, errors.New("refusing to compare non-stringer expected value")
	}
	return actualStringer.String() == expectedStringer.String(), nil
}

func (matcher *StringerEqualMatcher) FailureMessage(actual interface{}) (message string) {
	actualString, actualOK := actual.(string)
	expectedString, expectedOK := matcher.Expected.(string)
	if actualOK && expectedOK {
		return format.MessageWithDiff(actualString, "to equal", expectedString)
	}

	return format.Message(actual, "to equal", matcher.Expected)
}

func (matcher *StringerEqualMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to equal", matcher.Expected)
}
