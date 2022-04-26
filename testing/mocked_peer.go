package testing

import (
	"context"
	"fmt"
	"sync"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/msp"
	"github.com/pkg/errors"
)

type (
	MockedPeer struct {
		// channel name -> chaincode name
		ChannelCC ChannelsMockStubs
		m         sync.Mutex
	}

	ChannelMockStubs map[string]*MockStub

	ChannelsMockStubs map[string]ChannelMockStubs

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

func (mp *MockedPeer) WithChannel(channel string, mockStubs ...*MockStub) *MockedPeer {
	if _, ok := mp.ChannelCC[channel]; !ok {
		mp.ChannelCC[channel] = make(ChannelMockStubs)
	}

	for _, ms := range mockStubs {
		mp.ChannelCC[channel][ms.Name] = ms
		for chName, chnl := range mp.ChannelCC {
			for ccName, cc := range chnl {

				// add visibility of added cc to all other cc
				cc.MockPeerChaincode(ms.Name+`/`+channel, ms)

				// add visibility of other cc to added cc
				ms.MockPeerChaincode(ccName+`/`+chName, cc)
			}
		}
	}

	return mp
}

func (mp *MockedPeer) CurrentIdentity() msp.SigningIdentity {
	return nil
}

func (mp *MockedPeer) Invoke(
	ctx context.Context,
	channel string,
	chaincode string,
	args [][]byte,
	identity msp.SigningIdentity,
	transArgs map[string][]byte,
	txWaiterType string) (res *peer.Response, chaincodeTx string, err error) {

	mp.m.Lock()
	defer mp.m.Unlock()
	mockStub, err := mp.Chaincode(channel, chaincode)
	if err != nil {
		return nil, ``, err
	}

	response := mockStub.From(identity).WithTransient(transArgs).InvokeBytes(args...)
	if response.Status == shim.ERROR {
		err = errors.New(response.Message)
	}

	return &response, mockStub.TxID, err
}

func (mp *MockedPeer) Query(
	ctx context.Context,
	channel string,
	chaincode string,
	args [][]byte,
	identity msp.SigningIdentity,
	transArgs map[string][]byte) (*peer.Response, error) {

	mp.m.Lock()
	defer mp.m.Unlock()
	mockStub, err := mp.Chaincode(channel, chaincode)
	if err != nil {
		return nil, err
	}

	response := mockStub.From(identity).WithTransient(transArgs).QueryBytes(args...)
	if response.Status == shim.ERROR {
		err = errors.New(response.Message)
	}

	return &response, err

}

func (mp *MockedPeer) Chaincode(channel string, chaincode string) (*MockStub, error) {
	ms, exists := mp.ChannelCC[channel][chaincode]
	if !exists {
		return nil, fmt.Errorf(`%s: channell=%s, chaincode=%s`, ErrChaincodeNotExists, channel, chaincode)
	}

	return ms, nil
}

func (mp *MockedPeer) Events(
	ctx context.Context,
	channel string,
	chaincode string,
	identity msp.SigningIdentity,
	blockRange ...int64,
) (events chan interface {
	Event() *peer.ChaincodeEvent
	Block() uint64
	TxTimestamp() *timestamp.Timestamp
}, closer func() error, err error) {
	mockStub, err := mp.Chaincode(channel, chaincode)
	if err != nil {
		return nil, nil, err
	}

	var (
		eventsRaw chan *peer.ChaincodeEvent
	)
	// from oldest block to current channel height
	if len(blockRange) > 0 && blockRange[0] == 0 {
		// create copy of mockStub events chan
		eventsRaw, closer = mockStub.EventSubscription(0)

		// In real HLF peer stream not closed
		//// close events channel if its empty, because we receive events until current channel height
		//go func() {
		//	ticker := time.NewTicker(5 * time.Millisecond)
		//	for {
		//		<-ticker.C
		//		if len(events) == 0 {
		//			closer()
		//			ticker.Stop()
		//			return
		//		}
		//	}
		//}()
	} else {
		eventsRaw, closer = mockStub.EventSubscription()
	}

	go func() {
		<-ctx.Done()
		_ = closer()
	}()

	eventsExtended := make(chan interface {
		Event() *peer.ChaincodeEvent
		Block() uint64
		TxTimestamp() *timestamp.Timestamp
	})
	go func() {
		for e := range eventsRaw {
			eventsExtended <- &ChaincodeEvent{
				event: e,

				// todo: store information about block and timestamp in MockStub
				block:       55,
				txTimestamp: nil,
			}
		}
		close(eventsExtended)
	}()

	return eventsExtended, closer, nil
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
