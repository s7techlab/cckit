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

// EncMockStub wrapper for querying and invoking encrypted chaincode
type EncMockStub struct {
	CC     *testing.MockStub
	EncKey []byte
}

// NewEncMockStub creates wrapper for querying and invoking encrypted chaincode
func NewEncMockStub(cc *testing.MockStub, encKey []byte) *EncMockStub {
	return &EncMockStub{cc, encKey}
}

func (ecc *EncMockStub) Invoke(args ...interface{}) peer.Response {
	return MockInvoke(ecc.CC, ecc.EncKey, args...)
}

func (ecc *EncMockStub) Query(args ...interface{}) peer.Response {
	return MockQuery(ecc.CC, ecc.EncKey, args...)
}
