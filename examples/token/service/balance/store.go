package balance

import (
	"fmt"
	"strings"

	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state"
)

type Store struct {
	state state.GetSettable
}

func NewStore(ctx router.Context) *Store {
	return &Store{
		state: State(ctx),
	}
}

func (s *Store) Get(address string, token []string) (*Balance, error) {
	balance, err := s.state.Get(&BalanceId{Address: address, Token: token}, &Balance{})
	if err != nil {

		if strings.Contains(err.Error(), state.ErrKeyNotFound.Error()) {
			// default zero balance even if no Balance state entry for account exists
			return &Balance{
				Address: address,
				Amount:  0,
			}, nil
		}
		return nil, err
	}

	return balance.(*Balance), nil
}

func (s *Store) set(address string, token []string, amount uint64) error {
	balance := &Balance{
		Address: address,
		Token:   token,
		Amount:  amount,
	}

	return s.state.Put(balance)
}

func (s *Store) Add(address string, token []string, amount uint64) error {
	balance, err := s.Get(address, token)
	if err != nil {
		return nil
	}

	err = s.set(address, token, balance.Amount+amount)
	if err != nil {
		return fmt.Errorf(`add to=%s: %w`, address, err)
	}

	return nil
}

func (s *Store) Sub(address string, token []string, amount uint64) error {
	balance, err := s.Get(address, token)
	if err != nil {
		return err
	}

	if balance.Amount < amount {
		return fmt.Errorf(`subtract from=%s: %w`, address, ErrAmountInsuficcient)
	}

	err = s.set(address, token, balance.Amount-amount)
	if err != nil {
		return fmt.Errorf(`subtract from=%s: %w`, address, err)
	}

	return nil
}

func (s *Store) Transfer(senderAddress, recipientAddress string, tokenId []string, amount uint64) error {
	// subtract from sender balance
	if err := s.Sub(senderAddress, tokenId, amount); err != nil {
		return fmt.Errorf(`transfer: %w`, err)
	}
	// add to recipient balance
	if err := s.Add(recipientAddress, tokenId, amount); err != nil {
		return fmt.Errorf(`transfer: %w`, err)
	}
	return nil
}
