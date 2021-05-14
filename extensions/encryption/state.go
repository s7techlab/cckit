package encryption

import (
	"encoding/base64"
	"fmt"

	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/convert"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state"
)

var (
	// ErrKeyNotDefinedInTransientMap occurs when key not defined in transient map
	ErrKeyNotDefinedInTransientMap = errors.New(`encryption key is not defined in transient map`)
)

// State wrapper, encrypts the data before putting to state and
// decrypts the data after getting from state
func State(c router.Context, key []byte) (state.State, error) {
	//current state
	s := c.State()

	s.UseKeyTransformer(KeyEncryptor(key))
	s.UseKeyReverseTransformer(KeyDecryptor(key))
	s.UseStateGetTransformer(FromBytesDecryptor(key))
	s.UseStatePutTransformer(ToBytesEncryptor(key))

	return s, nil
}

// KeyFromTransient gets key for encrypting/decrypting from transient map
func KeyFromTransient(c router.Context) ([]byte, error) {
	tm, err := c.Stub().GetTransient()
	if err != nil {
		return nil, err
	}

	key, ok := tm[TransientMapKey]
	if !ok {
		return nil, ErrKeyNotDefinedInTransientMap
	}

	return key, nil
}

// StateWithTransientKey creates encrypted state state with provided key for symmetric encryption/decryption
func StateWithTransientKey(c router.Context) (state.State, error) {
	key, err := KeyFromTransient(c)
	if err != nil {
		return nil, err
	}
	return State(c, key)
}

// StateWithTransientKeyIfProvided creates encrypted state wrapper with provided key for symmetric encryption/decryption
// if key provided, otherwise - standard state wrapper without encryption
func StateWithTransientKeyIfProvided(c router.Context) (state.State, error) {
	key, err := KeyFromTransient(c)
	switch err {
	case nil:
		return State(c, key)
	case ErrKeyNotDefinedInTransientMap:
		//default state wrapper without encryption
		return c.State(), nil
	}
	return nil, err
}

// KeyEncryptor encrypts state key
func KeyEncryptor(encryptKey []byte) state.KeyTransformer {
	return func(key state.Key) (state.Key, error) {
		keyEnc := make(state.Key, len(key))

		for i, p := range key {
			keyPartEnc, err := Encrypt(encryptKey, p)
			if err != nil {
				return nil, fmt.Errorf(`encrypt key: %w`, err)
			}
			keyEnc[i] = base64.StdEncoding.EncodeToString(keyPartEnc)
		}
		return keyEnc, nil
	}
}

// KeyDecryptor decrypts state key
func KeyDecryptor(encryptKey []byte) state.KeyTransformer {
	return func(key state.Key) (state.Key, error) {
		keyEnc := make(state.Key, len(key))

		for i, p := range key {
			keyPartEnc, err := base64.StdEncoding.DecodeString(p)
			if err != nil {
				return nil, fmt.Errorf(`decrypt key base4 decode: %w`, err)
			}
			keyPart, err := Decrypt(encryptKey, keyPartEnc)
			if err != nil {
				return nil, fmt.Errorf(`decrypt key: %w`, err)
			}
			keyEnc[i] = string(keyPart)
		}
		return keyEnc, nil
	}
}

// DecryptTransformer returns state.FromBytesTransformer - used for decrypting data after reading from state
func FromBytesDecryptor(key []byte) state.FromBytesTransformer {
	return func(bb []byte, config ...interface{}) (interface{}, error) {
		decrypted, err := Decrypt(key, bb)
		if err != nil {
			return nil, errors.Wrap(err, `decrypt bytes`)
		}
		if len(config) == 0 {
			return decrypted, nil
		}
		return convert.FromBytes(decrypted, config[0])
	}
}

// EncryptTransformer returns state.ToBytesTransformer - used for encrypting data for state
func ToBytesEncryptor(key []byte) state.ToBytesTransformer {
	return func(v interface{}, config ...interface{}) ([]byte, error) {
		bb, err := convert.ToBytes(v)
		if err != nil {
			return nil, err
		}
		return Encrypt(key, bb)
	}
}

// EncryptWithTransientKey encrypts val with key from transient map
func EncryptWithTransientKey(c router.Context, val interface{}) (encrypted []byte, err error) {
	var (
		key []byte
	)

	if key, err = KeyFromTransient(c); err != nil {
		return
	}

	return ToBytesEncryptor(key)(val)
}
