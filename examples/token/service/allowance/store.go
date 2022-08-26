package allowance

import (
	"errors"
	"fmt"

	"github.com/s7techlab/cckit/examples/token/service/balance"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state"
)

type Store struct {
	state state.GetSettable
	path  []string // token path
}

func NewStore(ctx router.Context) *Store {
	return &Store{
		state: State(ctx),
	}
}

func (s *Store) Get(ownerAddress, spenderAddress string, token []string) (*Allowance, error) {
	allowance, err := s.state.Get(&AllowanceId{
		OwnerAddress:   ownerAddress,
		SpenderAddress: spenderAddress,
		Token:          token}, &Allowance{})
	if err != nil {
		if errors.Is(err, state.ErrKeyNotFound) {
			return &Allowance{
				OwnerAddress:   ownerAddress,
				SpenderAddress: spenderAddress,
				Token:          token,
				Amount:         0,
			}, nil
		}
		return nil, fmt.Errorf(`get allowance: %w`, err)
	}

	return allowance.(*Allowance), nil
}

func (s *Store) Set(ownerAddress, spenderAddress string, token []string, amount uint64) (*Allowance, error) {
	allowance := &Allowance{
		OwnerAddress:   ownerAddress,
		SpenderAddress: spenderAddress,
		Token:          token,
		Amount:         amount,
	}

	if err := s.state.Put(allowance); err != nil {
		return nil, fmt.Errorf(`set allowance: %w`, err)
	}

	return allowance, nil
}

func (s *Store) Sub(ownerAddress, spenderAddress string, token []string, amount uint64) (*Allowance, error) {
	allowance, err := s.Get(ownerAddress, spenderAddress, token)
	if err != nil {
		return nil, err
	}

	if allowance.Amount < amount {
		return nil, balance.ErrAmountInsufficient
	}

	return s.Set(ownerAddress, spenderAddress, token, allowance.Amount-amount)
}
