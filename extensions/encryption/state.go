package encryption

import (
	"github.com/s7techlab/cckit/convert"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state"
)

// State encrypting the data before putting to state and decrypting the data after getting from state
func State(c router.Context, key []byte) state.State {
	s := state.New(c.Stub())

	s.FromBytes = ConvertFromBytesWith(key)
	s.ToBytes = ConvertToBytesWith(key)

	return s
}

func EncryptState(state state.State) {

}

func ConvertFromBytesWith(key []byte) state.FromBytesTransformer {
	return state.ConvertFromBytes
}

func ConvertToBytesWith(key []byte) state.ToBytesTransformer {
	return func(v interface{}, config ...interface{}) ([]byte, error) {
		bb, err := convert.ToBytes(v)
		if err != nil {
			return nil, err
		}
		return Encrypt(key, bb)
	}
}
