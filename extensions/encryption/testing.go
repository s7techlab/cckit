package encryption

import (
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/cckit/response"
	"github.com/s7techlab/cckit/testing"
)

// MockInvoke helper for invoking MockStub with transient key and encrypted args
func MockInvoke(cc *testing.MockStub, encKey []byte, args ...interface{}) peer.Response {
	encArgs, err := EncryptArgs(encKey, args...)
	if err != nil {
		return response.Error(`unable to encrypt input args`)
	}
	return cc.WithTransient(TransientMapWithKey(encKey)).InvokeBytes(encArgs...)
}

// MockQuery helper for querying MockStub with transient key and encrypted args
func MockQuery(cc *testing.MockStub, encKey []byte, args ...interface{}) peer.Response {
	encArgs, err := EncryptArgs(encKey, args...)
	if err != nil {
		return response.Error(`unable to encrypt input args`)
	}
	return cc.WithTransient(TransientMapWithKey(encKey)).QueryBytes(encArgs...)
}

// MockStub wrapper for querying and invoking encrypted chaincode
type MockStub struct {
	MockStub *testing.MockStub
	EncKey   []byte
}

// NewMockStub creates wrapper for querying and invoking encrypted chaincode
func NewMockStub(mockStub *testing.MockStub, encKey []byte) *MockStub {
	return &MockStub{mockStub, encKey}
}

func (s *MockStub) Invoke(args ...interface{}) peer.Response {
	return MockInvoke(s.MockStub, s.EncKey, args...)
}

func (s *MockStub) Query(args ...interface{}) peer.Response {
	return MockQuery(s.MockStub, s.EncKey, args...)
}
