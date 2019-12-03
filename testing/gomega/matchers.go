package gomega

import (
	"github.com/onsi/gomega/types"
	"github.com/s7techlab/cckit/testing/gomega/matchers"
)

func StringerEqual(expected interface{}) types.GomegaMatcher {
	return &matchers.StringerEqualMatcher{
		Expected: expected,
	}
}
