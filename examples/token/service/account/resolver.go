package account

import (
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/s7techlab/cckit/router"
)

type Resolver interface {
	GetInvokerAddress(router.Context, *emptypb.Empty) (*AddressId, error)

	GetAddress(router.Context, *GetAddressRequest) (*AddressId, error)

	GetAccount(router.Context, *AccountId) (*Account, error)
}
