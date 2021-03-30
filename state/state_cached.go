package state

type (
	Cached struct {
		State
		TxCache map[string][]byte
	}
)

// WithCached returns state with tx level state cache
func WithCache(ss State) *Cached {
	s := ss.(*Impl)
	cached := &Cached{
		State:   s,
		TxCache: make(map[string][]byte),
	}

	s.PutState = func(key string, bb []byte) error {
		cached.TxCache[key] = bb
		return s.stub.PutState(key, bb)
	}

	s.GetState = func(key string) ([]byte, error) {
		if bb, ok := cached.TxCache[key]; ok {
			return bb, nil
		}
		return s.stub.GetState(key)
	}

	return cached
}
