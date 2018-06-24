package owner

import (
	"github.com/s7techlab/cckit/router"
)

const QueryMethod = `owner`

// FromState returns raw data ( serialized Grant ) of current chain code owner
func Query(c router.Context) (interface{}, error) {
	return c.State().Get(OwnerStateKey)
}
