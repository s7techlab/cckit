package encryption

import (
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/entities"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/convert"
)

const TransientMapKey = `ENCODE_KEY`

func init() {
	factory.InitFactories(nil)
}

func EncryptArgs(key []byte, args ...interface{}) ([][]byte, error) {
	argBytes, err := convert.ArgsToBytes(args...)
	if err != nil {
		return nil, err
	}

	for i, a := range argBytes {
		encrypted, err := Encrypt(key, a)
		if err != nil {
			return nil, errors.Wrap(err, `encryption error`)
		}
		argBytes[i] = encrypted
	}

	return argBytes, nil
}

func DecryptArgs(key []byte, args [][]byte) ([][]byte, error) {
	for i, a := range args {
		decrypted, err := Decrypt(key, a)
		if err != nil {
			return nil, errors.Wrap(err, `decryption error`)
		}
		args[i] = decrypted
	}

	return args, nil
}

func Encrypt(key, value []byte) ([]byte, error) {
	encrypter, err := entities.NewAES256EncrypterEntity("ID", factory.GetDefault(), key, nil)
	if err != nil {
		return nil, err
	}

	return encrypter.Encrypt(value)
}

func Decrypt(key, value []byte) ([]byte, error) {
	encrypter, err := entities.NewAES256EncrypterEntity("ID", factory.GetDefault(), key, nil)
	if err != nil {
		return nil, err
	}
	return encrypter.Decrypt(value)
}

func TransientMapWithKey(key []byte) map[string][]byte {
	return map[string][]byte{TransientMapKey: key}
}
