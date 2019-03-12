package encryption

import (
	"encoding/base64"

	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/entities"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/convert"
)

const TransientMapKey = `ENCODE_KEY`

func init() {
	factory.InitFactories(nil)
}

// EncryptArgs encrypt args
func EncryptArgs(key []byte, args ...interface{}) ([][]byte, error) {
	argBytes, err := convert.ArgsToBytes(args...)
	if err != nil {
		return nil, err
	}
	eargs := make([][]byte, len(args))
	for i, bb := range argBytes {
		encrypted, err := Encrypt(key, bb)
		if err != nil {
			return nil, errors.Wrap(err, `encryption error`)
		}

		eargs[i] = encrypted
	}
	return eargs, nil
}

// DecryptArgs decrypt args
func DecryptArgs(key []byte, args [][]byte) ([][]byte, error) {
	dargs := make([][]byte, len(args))
	for i, a := range args {
		decrypted, err := Decrypt(key, a)
		if err != nil {
			return nil, errors.Wrap(err, `decryption error`)
		}
		dargs[i] = decrypted
	}
	return dargs, nil
}

// Encrypt converts value to []byte  and encrypts its with key
func Encrypt(key []byte, value interface{}) ([]byte, error) {
	// TODO: customize  IV
	encrypter, err := entities.NewAES256EncrypterEntity("ID", factory.GetDefault(), key, make([]byte, 16))
	if err != nil {
		return nil, err
	}

	bb, err := convert.ToBytes(value)
	if err != nil {
		return nil, errors.Wrap(err, `convert values to bytes`)
	}

	return encrypter.Encrypt(bb)
}

// Decrypt decrypts value with key
func Decrypt(key, value []byte) ([]byte, error) {
	encrypter, err := entities.NewAES256EncrypterEntity("ID", factory.GetDefault(), key, nil)
	if err != nil {
		return nil, err
	}
	bb := make([]byte, len(value))
	copy(bb, value)
	return encrypter.Decrypt(bb)
}

// TransientMapWithKey creates transient map with encrypting/decrypting key
func TransientMapWithKey(key []byte) map[string][]byte {
	return map[string][]byte{TransientMapKey: key}
}

// EncryptEvent encrypts event payload and event name. Event name also base64 encoded.
// ChaincodeId and TxId remains unencrypted
func EncryptEvent(encKey []byte, event *peer.ChaincodeEvent) (encrypted *peer.ChaincodeEvent, err error) {
	var (
		encName, encPayload []byte
	)

	if encName, err = Encrypt(encKey, []byte(event.EventName)); err != nil {
		return nil, err
	}

	if encPayload, err = Encrypt(encKey, event.Payload); err != nil {
		return nil, err
	}

	return &peer.ChaincodeEvent{
		ChaincodeId: event.ChaincodeId,
		TxId:        event.TxId,
		EventName:   base64.StdEncoding.EncodeToString(encName),
		Payload:     encPayload,
	}, nil
}

// MustEncryptEvent helper for EncryptEvent. Panics in case of error.
func MustEncryptEvent(encKey []byte, event *peer.ChaincodeEvent) (encrypted *peer.ChaincodeEvent) {
	var (
		err error
	)
	if encrypted, err = EncryptEvent(encKey, event); err != nil {
		panic(err)
	}

	return encrypted
}
