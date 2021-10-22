package testing

import (
	"context"
	"fmt"
	"sync"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/pkg/errors"

	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/msp"
	"github.com/s7techlab/hlf-sdk-go/v2/api"
)

type (
	ChannelMockStubs map[string]*MockStub

	ChannelsMockStubs map[string]ChannelMockStubs

	MockedPeer struct {
		// channel name -> chaincode name
		ChannelCC ChannelsMockStubs
		m         sync.Mutex
	}

	EventSubscription struct {
		events chan *peer.ChaincodeEvent
		errors chan error
		closer sync.Once
	}
)

// NewInvoker implements Invoker interface from hlf-sdk-go
func NewPeer() *MockedPeer {
	return &MockedPeer{
		ChannelCC: make(ChannelsMockStubs),
	}
}

func (mi *MockedPeer) WithChannel(channel string, mockStubs ...*MockStub) *MockedPeer {
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

func (mi *MockedPeer) Invoke(
	ctx context.Context, from msp.SigningIdentity, channel string, chaincode string,
	fn string, args [][]byte, transArgs api.TransArgs, _ ...api.DoOption) (*peer.Response, api.ChaincodeTx, error) {

	mi.m.Lock()
	defer mi.m.Unlock()
	mockStub, err := mi.Chaincode(channel, chaincode)
	if err != nil {
		return nil, ``, err
	}

	response := mockStub.From(from).WithTransient(transArgs).InvokeBytes(append([][]byte{[]byte(fn)}, args...)...)
	if response.Status == shim.ERROR {
		err = errors.New(response.Message)
	}

	return &response, api.ChaincodeTx(mockStub.TxID), err
}

func (mi *MockedPeer) Query(
	ctx context.Context, from msp.SigningIdentity, channel string, chaincode string,
	fn string, args [][]byte, transArgs api.TransArgs) (*peer.Response, error) {
	mi.m.Lock()
	defer mi.m.Unlock()
	mockStub, err := mi.Chaincode(channel, chaincode)
	if err != nil {
		return nil, err
	}

	response := mockStub.From(from).WithTransient(transArgs).QueryBytes(append([][]byte{[]byte(fn)}, args...)...)
	if response.Status == shim.ERROR {
		err = errors.New(response.Message)
	}
	return &response, err
}

func (mi *MockedPeer) Subscribe(
	ctx context.Context, from msp.SigningIdentity, channel, chaincode string) (api.EventCCSubscription, error) {
	mockStub, err := mi.Chaincode(channel, chaincode)
	if err != nil {
		return nil, err
	}

	sub := &EventSubscription{
		events: mockStub.EventSubscription(),
		errors: make(chan error),
	}

	go func() {
		<-ctx.Done()
		close(sub.events)
		close(sub.errors)
	}()

	return sub, nil
}

func (mi *MockedPeer) Chaincode(channel string, chaincode string) (*MockStub, error) {
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
	es.closer.Do(func() {
		close(es.events)
		close(es.errors)
	})
	return nil
}
