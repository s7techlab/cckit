package mock

import (
	"github.com/s7techlab/cckit/gateway"
)

// Deprecated: use gateway.NewChaincodeService
func New(peer gateway.Peer) *gateway.ChaincodeService {
	return gateway.NewChaincodeService(peer)
}
