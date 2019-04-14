package encryption

import (
	"encoding/base64"

	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
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

// EventWithTransientKeyIfProvided returns encrypted event wrapper if key for symmetric
// encryption/decryption is provided, otherwise return default event wrapper
func EventWithTransientKeyIfProvided(c router.Context) (state.Event, error) {
	key, err := KeyFromTransient(c)
	switch err {
	case nil:
		return Event(c, key)
	case ErrKeyNotDefinedInTransientMap:
		//default event wrapper without encryption
		return c.Event(), nil
	}
	return nil, err
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

// DecryptEvent
func DecryptEvent(encKey []byte, event *peer.ChaincodeEvent) (decrypted *peer.ChaincodeEvent, err error) {
	var (
		encNameBytes, decName, decPayload []byte
	)

	if encNameBytes, err = base64.StdEncoding.DecodeString(event.EventName); err != nil {
		return nil, errors.Wrap(err, `event name base64 decoding`)
	}

	if decName, err = Decrypt(encKey, encNameBytes); err != nil {
		return nil, err
	}

	if decPayload, err = Decrypt(encKey, event.Payload); err != nil {
		return nil, err
	}

	return &peer.ChaincodeEvent{
		ChaincodeId: event.ChaincodeId,
		TxId:        event.TxId,
		EventName:   string(decName),
		Payload:     decPayload,
	}, nil
}
