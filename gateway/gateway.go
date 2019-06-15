package gateway

import (
	"context"

	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/cckit/extensions/encryption"
	"github.com/s7techlab/cckit/gateway/service"
)

// Chaincode interface for work with chaincode
type Chaincode interface {
	Query(ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error)
	Invoke(ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error)
	Events(ctx context.Context) (ChaincodeEventSub, error)
}

type ChaincodeEventSub interface {
	Context() context.Context
	Events() <-chan *peer.ChaincodeEvent
	Recv(*peer.ChaincodeEvent) error
	Close()
}

func WithDefaultSigner(defaultSigner msp.SigningIdentity) Opt {
	return func(c *chaincode) {
		c.ContextOpts = append(c.ContextOpts, func(ctx context.Context) context.Context {
			return service.ContextWithDefaultSigner(ctx, defaultSigner)
		})
	}
}

func WithTransientValue(key string, value []byte) Opt {
	return func(c *chaincode) {
		c.ContextOpts = append(c.ContextOpts, func(ctx context.Context) context.Context {
			return ContextWithTransientValue(ctx, key, value)
		})
	}
}

func WithEncryption(encKey []byte) Opt {
	return func(c *chaincode) {
		WithTransientValue(encryption.TransientMapKey, encKey)(c)
		WithArgsEncryption(encKey)(c)
		WithInvokePayloadDecryption(encKey)(c)
		WithEventDecryption(encKey)(c)
	}
}

func WithArgsEncryption(encKey []byte) Opt {
	return func(c *chaincode) {
		c.InputOpts = append(c.InputOpts, func(action Action, ccInput *service.ChaincodeInput) (err error) {
			ccInput.Args, err = encryption.EncryptArgsBytes(encKey, ccInput.Args)
			return err
		})
	}
}

func WithInvokePayloadDecryption(encKey []byte) Opt {
	return func(c *chaincode) {
		c.OutputOpts = append(c.OutputOpts, func(action Action, r *peer.Response) (err error) {
			if action != Invoke {
				return nil
			}
			r.Payload, err = encryption.Decrypt(encKey, r.Payload)
			return err
		})
	}
}

func WithEventDecryption(encKey []byte) Opt {
	return func(c *chaincode) {
		c.EventOpts = append(c.EventOpts, func(e *peer.ChaincodeEvent) error {
			de, err := encryption.DecryptEvent(encKey, e)
			if err != nil {
				return err
			}

			e.EventName = de.EventName
			e.Payload = de.Payload
			return nil
		})
	}
}
