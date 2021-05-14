package debug

import (
	"encoding/json"

	"github.com/golang/protobuf/proto"
	"github.com/s7techlab/cckit/convert"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state"
)

type (
	StateService struct {
		State StateFn
	}

	// StateFn function can add mappings to state, for correct convertation in StateGet
	StateFn func(router.Context) state.State
)

func StateAsIs(ctx router.Context) state.State {
	return ctx.State()
}

func NewStateService() *StateService {
	return &StateService{
		State: StateAsIs,
	}
}

func (s *StateService) StateClean(ctx router.Context, prefixes *Prefixes) (*PrefixesMatchCount, error) {
	var (
		keys []string
		key  string
		err  error
	)
	for _, p := range prefixes.Prefixes {
		if len(p.Key) > 0 {
			key, err = state.KeyToString(ctx.Stub(), p.Key)
			if err != nil {
				return nil, err
			}
		}

		keys = append(keys, key)
	}
	matches, err := DeleteStateByPrefixes(s.State(ctx), keys)
	if err != nil {
		return nil, err
	}

	return &PrefixesMatchCount{Matches: matches}, nil
}

func (s *StateService) StateKeys(ctx router.Context, prefix *Prefix) (*CompositeKeys, error) {
	keys, err := s.State(ctx).Keys(prefix.GetKey())
	if err != nil {
		return nil, err
	}

	cKeys := &CompositeKeys{}
	for _, keyStr := range keys {
		key, err := state.NormalizeKey(ctx.Stub(), keyStr)
		if err != nil {
			return nil, err
		}
		cKeys.Keys = append(cKeys.Keys, &CompositeKey{Key: key})
	}

	return cKeys, nil
}

func (s *StateService) StateGet(ctx router.Context, key *CompositeKey) (*Value, error) {
	val, err := s.State(ctx).Get(key.Key)
	if err != nil {
		return nil, err
	}

	bb, err := convert.ToBytes(val)
	if err != nil {
		return nil, err
	}

	var jsonVal []byte
	switch val.(type) {
	case proto.Message:
		jsonVal, _ = json.Marshal(val)
	}

	return &Value{
		Key:   key.Key,
		Value: bb,
		Json:  string(jsonVal),
	}, nil
}

func (s *StateService) StatePut(ctx router.Context, val *Value) (*Value, error) {
	if err := s.State(ctx).Put(val.Key, val.Value); err != nil {
		return nil, err
	}
	return val, nil
}

func (s *StateService) StateDelete(ctx router.Context, key *CompositeKey) (*Value, error) {
	val, err := s.StateGet(ctx, key)
	if err != nil {
		return nil, err
	}

	if err = s.State(ctx).Delete(key.Key); err != nil {
		return nil, err
	}

	return val, nil
}
