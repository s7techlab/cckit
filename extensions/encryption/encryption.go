package encryption

import (
	"fmt"

	"github.com/s7techlab/cckit/convert"
	"github.com/s7techlab/cckit/router"
)

const TransientMapKey = `ENCODE_KEY`

// EncryptArgs convert args to [][]byte and encrypt args with key
func EncryptArgs(key []byte, args ...interface{}) ([][]byte, error) {
	argBytes, err := convert.ArgsToBytes(args...)
	if err != nil {
		return nil, err
	}

	return EncryptArgsBytes(key, argBytes)
}

// EncryptArgsBytes encrypt args with key
func EncryptArgsBytes(key []byte, argsBytes [][]byte) ([][]byte, error) {
	eArgs := make([][]byte, len(argsBytes))
	for i, bb := range argsBytes {
		encrypted, err := EncryptBytes(key, bb)
		if err != nil {
			return nil, fmt.Errorf(`encryption error: %w`, err)
		}

		eArgs[i] = encrypted
	}
	return eArgs, nil
}

// DecryptArgs decrypt args
func DecryptArgs(key []byte, args [][]byte) ([][]byte, error) {
	dargs := make([][]byte, len(args))
	for i, a := range args {

		// do not try to decrypt init function
		if i == 0 && string(a) == router.InitFunc {
			dargs[i] = a
			continue
		}

		decrypted, err := DecryptBytes(key, a)
		if err != nil {
			return nil, fmt.Errorf(`decryption error: %w`, err)
		}
		dargs[i] = decrypted
	}
	return dargs, nil
}

// Encrypt converts value to []byte  and encrypts its with key
func Encrypt(key []byte, value interface{}) ([]byte, error) {
	// TODO: customize  IV
	bb, err := convert.ToBytes(value)
	if err != nil {
		return nil, fmt.Errorf(`convert values to bytes: %w`, err)
	}

	return EncryptBytes(key, bb)
}

// Decrypt decrypts value with key
func Decrypt(key, value []byte) ([]byte, error) {

	bb := make([]byte, len(value))
	copy(bb, value)
	return DecryptBytes(key, bb)
}

// TransientMapWithKey creates transient map with encrypting/decrypting key
func TransientMapWithKey(key []byte) map[string][]byte {
	return map[string][]byte{TransientMapKey: key}
}
