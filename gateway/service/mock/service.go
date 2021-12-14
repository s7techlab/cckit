package mock

import (
	"github.com/s7techlab/cckit/gateway/mock"
	"github.com/s7techlab/cckit/testing"
)

// Deprecated: use gateway/mock
func New(peers ...*testing.MockedPeer) *mock.ChaincodeService {
	return mock.New(peers...)
}
