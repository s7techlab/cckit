package allowance

import (
	"errors"
	"fmt"

	"github.com/s7techlab/cckit/examples/token/service/balance"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state"
	m "github.com/s7techlab/cckit/state/mapping"
)

var (
	StateMappings = m.StateMappings{}.
		//  Create mapping for Allowance entity
		// key `Allowance`,`{OwnerAddress}`,`{SpenderAddress}`,`{Path[0]}`..., `{Path[n]`
		Add(&Allowance{},
			m.PKeySchema(&AllowanceId{}),
			m.List(&Allowances{}), // Structure of result for List method
		)

	EventMappings = m.EventMappings{}.
			Add(&Approved{})
)

// State with chaincode mappings
func State(ctx router.Context) m.MappedState {
	return m.WrapState(ctx.State(), StateMappings)
}

// Event with chaincode mappings
func Event(ctx router.Context) state.Event {
	return m.WrapEvent(ctx.Event(), EventMappings)
}

type Store struct {
	state state.GetSettable
	path  []string // token path
}

func NewStore(ctx router.Context) *Store {
	return &Store{
		state: State(ctx),
	}
}

func (s *Store) Get(ownerAddress, spenderAddress string) (*Allowance, error) {
	allowance, err := s.state.Get(&AllowanceId{
		OwnerAddress:   ownerAddress,
		SpenderAddress: spenderAddress,
		Path:           s.path}, &Allowance{})
	if err != nil {
		if errors.Is(err, state.ErrKeyNotFound) {
			return &Allowance{
				OwnerAddress:   ownerAddress,
				SpenderAddress: spenderAddress,
				Amount:         0,
			}, nil
		}
		return nil, fmt.Errorf(`get allowance: %w`, err)
	}

	return allowance.(*Allowance), nil
}

func (s *Store) Set(ownerAddress, spenderAddress string, amount uint64) (*Allowance, error) {
	allowance := &Allowance{
		OwnerAddress:   ownerAddress,
		SpenderAddress: spenderAddress,
		Path:           s.path,
		Amount:         amount,
	}

	if err := s.state.Put(allowance); err != nil {
		return nil, fmt.Errorf(`set allowance: %w`, err)
	}

	return allowance, nil
}

func (s *Store) Sub(ownerAddress, spenderAddress string, amount uint64) (*Allowance, error) {
	allowance, err := s.Get(ownerAddress, spenderAddress)
	if err != nil {
		return nil, err
	}

	if allowance.Amount < amount {
		return nil, balance.ErrAmountInsuficcient
	}

	return s.Set(ownerAddress, spenderAddress, allowance.Amount-amount)
}
