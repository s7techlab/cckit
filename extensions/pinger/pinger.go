// Package pinger contains structure and functions for checking chain code accessibility
package pinger

import (
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/s7techlab/cckit/identity"
	r "github.com/s7techlab/cckit/router"
)

// Ping create PingInfo struct with tx creator ID and certificate in PEM format
func Ping(ctx r.Context) (interface{}, error) {
	id, err := cid.GetID(ctx.Stub())
	if err != nil {
		return nil, err
	}

	//take certificate from creator
	invoker, err := identity.FromStub(ctx.Stub())
	if err != nil {
		return nil, err
	}

	t, err := ctx.Time()
	if err != nil {
		return nil, err
	}

	return &PingInfo{
		InvokerId:   id,
		InvokerCert: invoker.GetPEM(),
		Time:        timestamppb.New(t),
	}, nil
}
