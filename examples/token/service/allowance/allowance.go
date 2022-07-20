package allowance

import (
	"errors"
	"fmt"

	"github.com/s7techlab/cckit/examples/token/service/balance"
	"github.com/s7techlab/cckit/router"
)

var (
	ErrOwnerOnly             = errors.New(`owner only`)
	ErrAllowanceInsufficient = errors.New(`allowance insufficient`)
)

type Service struct {
	balance *balance.Service
}

func NewService(balance *balance.Service) *Service {
	return &Service{
		balance: balance,
	}
}

func (s *Service) GetAllowance(ctx router.Context, req *AllowanceRequest) (*Allowance, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}

	allowance, err := NewStore(ctx).Get(req.OwnerAddress, req.SpenderAddress, req.Token)
	if err != nil {
		return nil, err
	}

	return allowance, nil
}

func (s *Service) Approve(ctx router.Context, req *ApproveRequest) (*Allowance, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}

	invokerAddress, err := s.balance.Account.GetInvokerAddress(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf(`get invoker address: %w`, err)
	}

	if invokerAddress.Address != req.OwnerAddress {
		return nil, ErrOwnerOnly
	}

	allowance, err := NewStore(ctx).Set(req.OwnerAddress, req.SpenderAddress, req.Token, req.Amount)
	if err != nil {
		return nil, err
	}

	if err = Event(ctx).Set(&Approved{
		OwnerAddress:   req.OwnerAddress,
		SpenderAddress: req.SpenderAddress,
		Amount:         req.Amount,
	}); err != nil {
		return nil, err
	}

	return allowance, nil
}

func (s *Service) TransferFrom(ctx router.Context, req *TransferFromRequest) (*TransferFromResponse, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}

	spenderAddress, err := s.balance.Account.GetInvokerAddress(ctx, nil)
	if err != nil {
		return nil, err
	}

	allowance, err := NewStore(ctx).Get(req.OwnerAddress, spenderAddress.Address, req.Token)
	if err != nil {
		return nil, err
	}

	if allowance.Amount < req.Amount {
		return nil, fmt.Errorf(`request trasfer amount=%d, allowance=%d: %w`,
			req.Amount, allowance.Amount, ErrAllowanceInsufficient)
	}

	if err = s.balance.Store(ctx).Transfer(req.OwnerAddress, req.RecipientAddress, req.Token, req.Amount); err != nil {
		return nil, err
	}

	if err = Event(ctx).Set(&TransferredFrom{
		OwnerAddress:     req.OwnerAddress,
		SpenderAddress:   spenderAddress.Address,
		RecipientAddress: req.RecipientAddress,
		Amount:           0,
	}); err != nil {
		return nil, err
	}

	return &TransferFromResponse{
		OwnerAddress:     req.OwnerAddress,
		RecipientAddress: req.RecipientAddress,
		Amount:           req.Amount,
	}, nil
}
