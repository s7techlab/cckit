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
	//EncKey key for encrypt data before query/invoke
	EncKey []byte

	// DecryptInvokeResponse decrypts invoker responses
	DecryptInvokeResponse bool
}

// NewMockStub creates wrapper for querying and invoking encrypted chaincode
func NewMockStub(mockStub *testing.MockStub, encKey []byte) *MockStub {
	return &MockStub{MockStub: mockStub, EncKey: encKey}
}

func (s *MockStub) Invoke(args ...interface{}) (response peer.Response) {
	var (
		err error
	)
	// first we encrypt all args
	response = MockInvoke(s.MockStub, s.EncKey, args...)

	//after receiving response we can decrypt received peer responce
	// actual only for invoke, query responses are not encrypted
	if s.DecryptInvokeResponse && len(response.Payload) > 0 {
		if response.Payload, err = Decrypt(s.EncKey, response.Payload); err != nil {
			panic(err)
		}
	}

	return response
}

func (s *MockStub) Query(args ...interface{}) peer.Response {
	return MockQuery(s.MockStub, s.EncKey, args...)
}

func (s *MockStub) From(args ...interface{}) *MockStub {
	s.MockStub.From(args...)
	return s
}
