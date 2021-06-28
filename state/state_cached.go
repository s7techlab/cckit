package state

type (
	Cached struct {
		State
		TxWriteSet  map[string][]byte
		TxDeleteSet map[string]interface{}
	}
)

// WithCached returns state with tx level state cache
func WithCache(ss State) *Cached {
	s := ss.(*Impl)
	cached := &Cached{
		State:       s,
		TxWriteSet:  make(map[string][]byte),
		TxDeleteSet: make(map[string]interface{}),
	}

	s.PutState = func(key string, bb []byte) error {
		cached.TxWriteSet[key] = bb
		return s.stub.PutState(key, bb)
	}

	s.GetState = func(key string) ([]byte, error) {
		if bb, ok := cached.TxWriteSet[key]; ok {
			return bb, nil
		}

		if _, ok := cached.TxDeleteSet[key]; ok {
			return []byte{}, nil
		}
		return s.stub.GetState(key)
	}

	s.DeleteState = func(key string) error {
		delete(cached.TxWriteSet, key)
		cached.TxDeleteSet[key] = nil

		return s.stub.DelState(key)
	}

	return cached
}
