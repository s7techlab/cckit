package testing

import (
	"fmt"

	"github.com/hyperledger/fabric-protos-go/peer"

	"github.com/s7techlab/cckit/extensions/encryption"
	"github.com/s7techlab/cckit/response"
	testcc "github.com/s7techlab/cckit/testing"
)

// MockInvoke helper for invoking MockStub with transient key and encrypted args
func MockInvoke(cc *testcc.MockStub, encKey []byte, args ...interface{}) peer.Response {
	encArgs, err := encryption.EncryptArgs(encKey, args...)
	if err != nil {
		return response.Error(`unable to encrypt input args`)
	}
	return cc.AddTransient(encryption.TransientMapWithKey(encKey)).InvokeBytes(encArgs...)
}

// MockQuery helper for querying MockStub with transient key and encrypted args
func MockQuery(cc *testcc.MockStub, encKey []byte, args ...interface{}) peer.Response {
	encArgs, err := encryption.EncryptArgs(encKey, args...)
	if err != nil {
		return response.Error(`unable to encrypt input args`)
	}
	return cc.AddTransient(encryption.TransientMapWithKey(encKey)).QueryBytes(encArgs...)
}

// MockStub wrapper for querying and invoking encrypted chaincode
type MockStub struct {
	MockStub *testcc.MockStub
	//EncKey key for encrypt data before query/invoke
	EncKey []byte

	// DecryptInvokeResponse decrypts invoker responses
	DecryptInvokeResponse bool
}

// NewMockStub creates wrapper for querying and invoking encrypted chaincode
func NewMockStub(mockStub *testcc.MockStub, encKey []byte) *MockStub {
	return &MockStub{MockStub: mockStub, EncKey: encKey}
}

func (s *MockStub) Invoke(args ...interface{}) (response peer.Response) {
	var (
		err       error
		decrypted []byte
	)
	// first we encrypt all args
	response = MockInvoke(s.MockStub, s.EncKey, args...)

	//after receiving response we can decrypt received peer response
	// actual only for invoke, query responses are not encrypted
	if s.DecryptInvokeResponse && len(response.Payload) > 0 && string(response.Payload) != `null` {
		if decrypted, err = encryption.Decrypt(s.EncKey, response.Payload); err != nil {
			panic(fmt.Sprintf(
				`decrypt mock invoke error with payload %s (%d): %s`,
				string(response.Payload), len(response.Payload), err))
		}
		response.Payload = decrypted
	}

	return response
}

func (s *MockStub) Query(args ...interface{}) peer.Response {
	return MockQuery(s.MockStub, s.EncKey, args...)
}

func (s *MockStub) Init(args ...interface{}) peer.Response {
	encArgs, err := encryption.EncryptArgs(s.EncKey, args...)
	if err != nil {
		return response.Error(`unable to encrypt input args`)
	}
	return s.MockStub.AddTransient(encryption.TransientMapWithKey(s.EncKey)).InitBytes(encArgs...)
}

func (s *MockStub) From(args ...interface{}) *MockStub {
	s.MockStub.From(args...)
	return s
}

func (s *MockStub) LastEvent() *peer.ChaincodeEvent {
	return encryption.MustDecryptEvent(s.EncKey, s.MockStub.ChaincodeEvent)
}
