package testing

import (
	"context"
	"fmt"
	"sync"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/msp"
	"github.com/pkg/errors"
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

// NewPeer implements Peer interface
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
	ctx context.Context,
	channel string,
	chaincode string,
	args [][]byte,
	identity msp.SigningIdentity,
	transArgs map[string][]byte,
	txWaiterType string) (res *peer.Response, chaincodeTx string, err error) {

	mi.m.Lock()
	defer mi.m.Unlock()
	mockStub, err := mi.Chaincode(channel, chaincode)
	if err != nil {
		return nil, ``, err
	}

	response := mockStub.From(identity).WithTransient(transArgs).InvokeBytes(args...)
	if response.Status == shim.ERROR {
		err = errors.New(response.Message)
	}

	return &response, mockStub.TxID, err
}

func (mi *MockedPeer) Query(
	ctx context.Context,
	channel string,
	chaincode string,
	args [][]byte,
	identity msp.SigningIdentity,
	transArgs map[string][]byte) (*peer.Response, error) {

	mi.m.Lock()
	defer mi.m.Unlock()
	mockStub, err := mi.Chaincode(channel, chaincode)
	if err != nil {
		return nil, err
	}

	response := mockStub.From(identity).WithTransient(transArgs).QueryBytes(args...)
	if response.Status == shim.ERROR {
		err = errors.New(response.Message)
	}

	return &response, err

}

func (mi *MockedPeer) Chaincode(channel string, chaincode string) (*MockStub, error) {
	ms, exists := mi.ChannelCC[channel][chaincode]
	if !exists {
		return nil, fmt.Errorf(`%s: channell=%s, chaincode=%s`, ErrChaincodeNotExists, channel, chaincode)
	}

	return ms, nil
}

func (mi *MockedPeer) Events(
	ctx context.Context,
	channel string,
	chaincode string,
	identity msp.SigningIdentity,
	blockRange ...int64,
) (chan *peer.ChaincodeEvent, error) {

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

	return sub.Events(), nil
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
