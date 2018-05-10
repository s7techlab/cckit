package owner

import (
	"github.com/s7techlab/cckit/router"
)

// FromState returns raw data ( serialized Grant ) of current chain code owner
func FromState(c router.Context) (interface{}, error) {
	return c.State().Get(OwnerStateKey)
}
