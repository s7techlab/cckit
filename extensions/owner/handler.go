package owner

import (
	"github.com/s7techlab/cckit/router"
)

const QueryMethod = `owner`

// Query returns raw data (serialized Grant) of current chain code owner
func Query(c router.Context) (interface{}, error) {
	return c.State().Get(OwnerStateKey)
}

// InvokeSetFromCreator sets tx creator as chaincode owner, if owner not previously set
func InvokeSetFromCreator(c router.Context) (interface{}, error) {
	return SetFromCreator(c)
}

// InvokeSetFromArgs gets owner data from args[0] (Msp Id) and arg[1] (cert)
func InvokeSetFromArgs(c router.Context) (interface{}, error) {
	return SetFromArgs(c)
}
