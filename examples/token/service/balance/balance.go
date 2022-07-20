package balance

import (
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/s7techlab/cckit/examples/token/service/account"
	"github.com/s7techlab/cckit/examples/token/service/config"

	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state"
)

type Service struct {
	Account account.Getter
	Token   config.TokenGetter
}

func New(accountResolver account.Getter, tokenGetter config.TokenGetter) *Service {
	return &Service{
		Account: accountResolver,
		Token:   tokenGetter,
	}
}

func (s *Service) Store(ctx router.Context) *Store {
	return NewStore(ctx)
}

func (s *Service) GetBalance(ctx router.Context, req *GetBalanceRequest) (*Balance, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}

	token, err := s.Token.GetToken(ctx, &config.TokenId{Token: req.Token})
	if err != nil {
		return nil, fmt.Errorf(`get token: %w`, err)
	}
	return s.Store(ctx).Get(req.Address, token.Token)
}

func (s *Service) ListBalances(ctx router.Context, _ *emptypb.Empty) (*Balances, error) {
	balances, err := State(ctx).List(&Balance{})
	if err != nil {
		return nil, err
	}
	return balances.(*Balances), nil
}

func (s *Service) ListAddressBalances(ctx router.Context, req *ListAddressBalancesRequest) (*Balances, error) {
	balances, err := State(ctx).ListWith(&Balance{}, state.Key{req.Address})
	if err != nil {
		return nil, err
	}
	return balances.(*Balances), nil
}

func (s *Service) Transfer(ctx router.Context, req *TransferRequest) (*TransferResponse, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}

	invokerAddress, err := s.Account.GetInvokerAddress(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf(`get invoker address: %w`, err)
	}

	token, err := s.Token.GetToken(ctx, &config.TokenId{Token: req.Token})
	if err != nil {
		return nil, fmt.Errorf(`get token: %w`, err)
	}

	if err := s.Store(ctx).Transfer(invokerAddress.Address, req.RecipientAddress, token.Token, req.Amount); err != nil {
		return nil, err
	}

	if err = Event(ctx).Set(&Transferred{
		SenderAddress:    invokerAddress.Address,
		RecipientAddress: req.RecipientAddress,
		Token:            token.Token,
		Amount:           req.Amount,
	}); err != nil {
		return nil, err
	}

	return &TransferResponse{
		SenderAddress:    invokerAddress.Address,
		RecipientAddress: req.RecipientAddress,
		Token:            token.Token,
		Amount:           req.Amount,
	}, nil
}
