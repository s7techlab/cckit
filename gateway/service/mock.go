package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/cckit/testing"
)

const (
	MessageProtocolVersion = 1
)

type (
	ChannelMockStubs map[string]*testing.MockStub

	ChannelsMockStubs map[string]ChannelMockStubs

	MockChaincodeService struct {
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

func NewMock() *MockChaincodeService {
	return &MockChaincodeService{
		ChannelCC: make(ChannelsMockStubs),
	}
}
func (cs *MockChaincodeService) Query(ctx context.Context, in *ChaincodeInput) (proposalResponse *peer.ProposalResponse, err error) {
	var (
		mockStub *testing.MockStub
		signer   msp.SigningIdentity
		response peer.Response
	)

	cs.m.Lock()
	defer cs.m.Unlock()

	if mockStub, err = cs.Chaincode(in.Channel, in.Chaincode); err != nil {
		return
	}

	if signer, err = SignerFromContext(ctx); err != nil {
		return
	}

	if response = mockStub.From(signer).WithTransient(in.Transient).QueryBytes(in.Args...); response.Status == shim.ERROR {
		return nil, errors.New(response.Message)
	}

	return &peer.ProposalResponse{
		Version:   MessageProtocolVersion,
		Timestamp: mockStub.TxTimestamp,
		Response:  &response,
	}, nil
}

func (cs *MockChaincodeService) Invoke(ctx context.Context, in *ChaincodeInput) (proposalResponse *peer.ProposalResponse, err error) {
	var (
		mockStub *testing.MockStub
		signer   msp.SigningIdentity
		response peer.Response
	)

	cs.m.Lock()
	defer cs.m.Unlock()

	if mockStub, err = cs.Chaincode(in.Channel, in.Chaincode); err != nil {
		return
	}

	if signer, err = SignerFromContext(ctx); err != nil {
		return
	}

	if response = mockStub.From(signer).WithTransient(in.Transient).InvokeBytes(in.Args...); response.Status == shim.ERROR {
		return nil, errors.New(response.Message)
	}

	return &peer.ProposalResponse{
		Version:   MessageProtocolVersion,
		Timestamp: mockStub.TxTimestamp,
		Response:  &response,
	}, nil

}

func (cs *MockChaincodeService) Events(in *ChaincodeLocator, stream Chaincode_EventsServer) (err error) {
	var (
		mockStub *testing.MockStub
	)
	if mockStub, err = cs.Chaincode(in.Channel, in.Chaincode); err != nil {
		return
	}
	events := mockStub.EventSubscription()
	for {
		e, ok := <-events
		if !ok {
			return nil
		}
		if err = stream.Send(e); err != nil {
			return err
		}
	}
}

func (cs *MockChaincodeService) WithChannel(channel string, mockStubs ...*testing.MockStub) *MockChaincodeService {
	if _, ok := cs.ChannelCC[channel]; !ok {
		cs.ChannelCC[channel] = make(ChannelMockStubs)
	}

	for _, ms := range mockStubs {
		cs.ChannelCC[channel][ms.Name] = ms
		for chName, chnl := range cs.ChannelCC {
			for ccName, cc := range chnl {

				// add visibility of added cc to all other cc
				cc.MockPeerChaincode(ms.Name+`/`+channel, ms)

				// add visibility of other cc to added cc
				ms.MockPeerChaincode(ccName+`/`+chName, cc)
			}
		}
	}

	return cs
}

func (cs *MockChaincodeService) Chaincode(channel string, chaincode string) (*testing.MockStub, error) {
	ms, exists := cs.ChannelCC[channel][chaincode]
	if !exists {
		return nil, fmt.Errorf(`%s: channell=%s, chaincode=%s`, ErrChaincodeNotExists, channel, chaincode)
	}

	return ms, nil
}
