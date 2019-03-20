package testing

import (
	"context"
	"fmt"

	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/protos/peer"
)

type (
	// Types from hlf-sdk-go

	ChaincodeTx string

	// Invoker interface describes common operations for chaincode
	Invoker interface {
		// Invoke method allows to invoke chaincode
		Invoke(ctx context.Context, from msp.SigningIdentity, channel string, chaincode string, fn string, args [][]byte) (*peer.Response, ChaincodeTx, error)
		// Query method allows to query chaincode without sending response to orderer
		Query(ctx context.Context, from msp.SigningIdentity, channel string, chaincode string, fn string, args [][]byte) (*peer.Response, error)
		// Subscribe allows to subscribe on chaincode events
		Subscribe(ctx context.Context, from msp.SigningIdentity, channel, chaincode string) (EventCCSubscription, error)
	}

	// EventCCSubscription describes chaincode events subscription
	EventCCSubscription interface {
		// Events initiates internal GRPC stream and returns channel on chaincode events
		Events() chan *peer.ChaincodeEvent
		// Errors returns errors associated with this subscription
		Errors() chan error
		// Close cancels current subscription
		Close() error
	}

	ChannelMockStubs map[string]*MockStub

	ChannelsMockStubs map[string]ChannelMockStubs

	MockInvoker struct {
		// channel name -> chaincode name
		ChannelCC ChannelsMockStubs
	}

	EventSubscription struct {
		events chan *peer.ChaincodeEvent
		errors chan error
	}
)

// NewInvoker implements hlf-sdk-go
func NewInvoker() *MockInvoker {
	return &MockInvoker{
		ChannelCC: make(ChannelsMockStubs),
	}
}

func (mi *MockInvoker) WithChannel(channel string, mockStubs ...*MockStub) *MockInvoker {
	if _, ok := mi.ChannelCC[channel]; !ok {
		mi.ChannelCC[channel] = make(ChannelMockStubs)
	}

	for _, ms := range mockStubs {
		mi.ChannelCC[channel][ms.Name] = ms
		for chName, chnl := range mi.ChannelCC {
			for ccName, cc := range chnl {

				// add visibility of added cc to all other cc
				cc.MockPeerChaincode(ms.Name+`/`+channel, ms)

				// add visibility of other cc to added cc
				ms.MockPeerChaincode(ccName+`/`+chName, cc)
			}
		}
	}

	return mi
}

func (mi *MockInvoker) Invoke(ctx context.Context, from msp.SigningIdentity, channel string, chaincode string, fn string, args [][]byte) (*peer.Response, ChaincodeTx, error) {
	mockStub, err := mi.Chaincode(channel, chaincode)
	if err != nil {
		return nil, ``, err
	}

	response := mockStub.From(from).InvokeBytes(append([][]byte{[]byte(fn)}, args...)...)

	return &response, ChaincodeTx(mockStub.TxID), nil
}

func (mi *MockInvoker) Query(ctx context.Context, from msp.SigningIdentity, channel string, chaincode string, fn string, args [][]byte) (*peer.Response, error) {
	mockStub, err := mi.Chaincode(channel, chaincode)
	if err != nil {
		return nil, err
	}

	response := mockStub.From(from).QueryBytes(append([][]byte{[]byte(fn)}, args...)...)
	return &response, nil
}

func (mi *MockInvoker) Subscribe(ctx context.Context, from msp.SigningIdentity, channel, chaincode string) (EventCCSubscription, error) {
	mockStub, err := mi.Chaincode(channel, chaincode)
	if err != nil {
		return nil, err
	}

	return &EventSubscription{
		events: mockStub.EventSubscription(),
		errors: make(chan error),
	}, nil
}

func (mi *MockInvoker) Chaincode(channel string, chaincode string) (*MockStub, error) {
	ms, exists := mi.ChannelCC[channel][chaincode]
	if !exists {
		return nil, fmt.Errorf(`%s: channell=%s, chaincode=%s`, ErrChaincodeNotExists, channel, chaincode)
	}

	return ms, nil
}

func (es *EventSubscription) Events() chan *peer.ChaincodeEvent {
	return es.events
}

func (es *EventSubscription) Errors() chan error {
	return es.errors
}

func (es *EventSubscription) Close() error {
	return nil
}
