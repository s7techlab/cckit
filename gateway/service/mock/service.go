package mock

import (
	"github.com/s7techlab/cckit/gateway/mock"
)

// Deprecated: use gateway/mock
func New() *mock.ChaincodeService {
	return mock.New()
}
