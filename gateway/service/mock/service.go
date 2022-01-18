package mock

import (
	"github.com/s7techlab/cckit/gateway"
	"github.com/s7techlab/cckit/sdk"
)

// Deprecated: use gateway.NewChaincodeService
func New(sdk sdk.SDK) *gateway.ChaincodeService {
	return gateway.NewChaincodeService(sdk)
}
