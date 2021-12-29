package encryption

import (
	"github.com/hyperledger/fabric-protos-go/peer"
)

// MustEncryptEvent helper for EncryptEvent. Panics in case of error.
func MustEncryptEvent(encKey []byte, event *peer.ChaincodeEvent) *peer.ChaincodeEvent {
	encrypted, err := EncryptEvent(encKey, event)
	if err != nil {
		panic(err)
	}
	return encrypted
}

// MustDecryptEvent helper for DecryptEvent. Panics in case of error.
func MustDecryptEvent(encKey []byte, event *peer.ChaincodeEvent) *peer.ChaincodeEvent {
	decrypted, err := DecryptEvent(encKey, event)
	if err != nil {
		panic(err)
	}

	return decrypted
}
