package encryption

import (
	"encoding/base64"

	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state"
)

// EventWithTransientKey creates encrypted event wrapper with provided key for symmetric encryption/decryption
func EventWithTransientKey(c router.Context) (state.Event, error) {
	key, err := KeyFromTransient(c)
	if err != nil {
		return nil, err
	}
	return Event(c, key)
}

// Event encrypting the events before setEvent()
func Event(c router.Context, key []byte) (state.Event, error) {
	//current state
	s := c.Event()
	s.UseSetTransformer(ToBytesEncryptor(key))
	s.UseNameTransformer(StringEncryptor(key))
	return s, nil
}

// EncryptStringWith returns state.StringTransformer encrypting string with provided key
func StringEncryptor(key []byte) state.StringTransformer {
	return func(s string) (encrypted string, err error) {
		var (
			enc []byte
		)
		if enc, err = Encrypt(key, []byte(s)); err != nil {
			return ``, err
		}

		return base64.StdEncoding.EncodeToString(enc), nil
	}
}
