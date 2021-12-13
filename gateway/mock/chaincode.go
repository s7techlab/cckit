package mock

import (
	"context"
	"errors"
	"sync"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/msp"
	"github.com/s7techlab/cckit/gateway"

	"github.com/s7techlab/cckit/testing"
)

const (
	MessageProtocolVersion = 1
)

type (
	ChaincodeService struct {
		gateway.ChaincodeService
		// channel name -> chaincode name
		Peer    *testing.MockedPeer
		m       sync.Mutex
		Invoker ChaincodeInvoker
	}

	// ChaincodeInvoker allows to imitate peer errors or unavailability
	ChaincodeInvoker func(ctx context.Context, mockStub *testing.MockStub, in *gateway.ChaincodeExec) *peer.Response
)

func New(peers ...*testing.MockedPeer) *ChaincodeService {
	var p *testing.MockedPeer
	if len(peers) > 0 {
		p = peers[0]
	} else {
		p = testing.NewPeer()
	}

	return &ChaincodeService{
		Peer:    p,
		Invoker: DefaultInvoker,
	}
}

func DefaultInvoker(ctx context.Context, mockStub *testing.MockStub, in *gateway.ChaincodeExec) *peer.Response {
	var (
		signer   msp.SigningIdentity
		response peer.Response
		err      error
	)

	if signer, err = gateway.SignerFromContext(ctx); err != nil {
		return &peer.Response{
			Status:  shim.ERROR,
			Message: `signer is not defined in context`,
		}
	}

	mockStub.From(signer).WithTransient(in.Input.Transient)

	if in.Type == gateway.InvocationType_QUERY {
		response = mockStub.QueryBytes(in.Input.Args...)
	} else if in.Type == gateway.InvocationType_INVOKE {
		response = mockStub.InvokeBytes(in.Input.Args...)
	} else {
		return &peer.Response{
			Status:  shim.ERROR,
			Message: gateway.ErrUnknownInvocationType.Error(),
		}
	}

	return &response
}

func (cs *ChaincodeService) Exec(ctx context.Context, in *gateway.ChaincodeExec) (*peer.Response, error) {
	var (
		mockStub *testing.MockStub
		response *peer.Response
		err      error
	)

	cs.m.Lock()
	defer cs.m.Unlock()

	if mockStub, err = cs.Peer.Chaincode(in.Input.Chaincode.Channel, in.Input.Chaincode.Chaincode); err != nil {
		return nil, err
	}

	response = cs.Invoker(ctx, mockStub, in)

	if response.Status == shim.ERROR {
		return nil, errors.New(response.Message)
	}

	return response, nil
}

func (cs *ChaincodeService) Query(ctx context.Context, in *gateway.ChaincodeInput) (*peer.Response, error) {
	return cs.Exec(ctx, &gateway.ChaincodeExec{
		Type:  gateway.InvocationType_QUERY,
		Input: in,
	})
}

func (cs *ChaincodeService) Invoke(ctx context.Context, in *gateway.ChaincodeInput) (proposalResponse *peer.Response, err error) {
	return cs.Exec(ctx, &gateway.ChaincodeExec{
		Type:  gateway.InvocationType_INVOKE,
		Input: in,
	})
}

func (cs *ChaincodeService) Events(in *gateway.ChaincodeEventsRequest, stream gateway.ChaincodeService_EventsServer) (err error) {
	var (
		mockStub *testing.MockStub
	)
	if mockStub, err = cs.Peer.Chaincode(in.Chaincode.Channel, in.Chaincode.Chaincode); err != nil {
		return
	}
	events := mockStub.EventSubscription()
	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case e, ok := <-events:
			if !ok {
				return nil
			}
			if err = stream.Send(e); err != nil {
				return err
			}
		}
	}
}
