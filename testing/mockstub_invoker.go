package testing

import (
	"context"
	"fmt"
	"sync"

	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/hlf-sdk-go/api"
)

type (
	ChannelMockStubs map[string]*MockStub

	ChannelsMockStubs map[string]ChannelMockStubs

	MockInvoker struct {
		// channel name -> chaincode name
		ChannelCC ChannelsMockStubs
		m         sync.Mutex
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

func (mi *MockInvoker) Invoke(ctx context.Context, from msp.SigningIdentity, channel string, chaincode string, fn string, args [][]byte, transArgs api.TransArgs) (*peer.Response, api.ChaincodeTx, error) {
	mi.m.Lock()
	defer mi.m.Unlock()
	mockStub, err := mi.Chaincode(channel, chaincode)
	if err != nil {
		return nil, ``, err
	}

	response := mockStub.From(from).WithTransient(transArgs).InvokeBytes(append([][]byte{[]byte(fn)}, args...)...)
	return &response, api.ChaincodeTx(mockStub.TxID), nil
}

func (mi *MockInvoker) Query(ctx context.Context, from msp.SigningIdentity, channel string, chaincode string, fn string, args [][]byte, transArgs api.TransArgs) (*peer.Response, error) {
	mi.m.Lock()
	defer mi.m.Unlock()
	mockStub, err := mi.Chaincode(channel, chaincode)
	if err != nil {
		return nil, err
	}

	response := mockStub.From(from).WithTransient(transArgs).QueryBytes(append([][]byte{[]byte(fn)}, args...)...)
	return &response, nil
}

func (mi *MockInvoker) Subscribe(ctx context.Context, from msp.SigningIdentity, channel, chaincode string) (api.EventCCSubscription, error) {
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
