package owner

import (
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/cckit/extensions/access"
	"github.com/s7techlab/cckit/router"
)

// Get chaincode invoke handler returns current chain code owner
func Get(c router.Context) peer.Response {
	return c.Response().Create(c.State().Get(OwnerStateKey, access.Grant{}))
}
