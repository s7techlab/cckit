package account

import (
	"encoding/base64"
	"errors"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/s7techlab/cckit/identity"
	"github.com/s7techlab/cckit/router"
)

type LocalService struct {
}

func NewLocalService() *LocalService {
	return &LocalService{}
}

func (l *LocalService) GetInvokerAddress(ctx router.Context, _ *emptypb.Empty) (*AddressId, error) {
	invoker, err := identity.FromStub(ctx.Stub())
	if err != nil {
		return nil, err
	}

	return l.GetAddress(ctx,
		&GetAddressRequest{PublicKey: identity.MarshalPublicKey(invoker.Cert.PublicKey)})
}

func (l *LocalService) GetAddress(ctx router.Context, req *GetAddressRequest) (*AddressId, error) {
	return &AddressId{
		Address: base64.StdEncoding.EncodeToString(req.PublicKey),
	}, nil
}

func (l *LocalService) GetAccount(ctx router.Context, id *AccountId) (*Account, error) {
	return nil, errors.New(`no accounts implemented`)
}

type Remote struct {
}
