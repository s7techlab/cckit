package gateway

import (
	"context"

	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/msp"

	"github.com/s7techlab/cckit/convert"
	"github.com/s7techlab/cckit/extensions/encryption"
	"github.com/s7techlab/cckit/state/mapping"
)

type (
	Opts struct {
		Context []ContextOpt
		Input   []InputOpt
		Output  []OutputOpt
		Event   []EventOpt
	}

	Opt func(*Opts)

	ContextOpt func(ctx context.Context) context.Context
	InputOpt   func(action Action, input *ChaincodeInput) error
	OutputOpt  func(action Action, response *peer.Response) error
	EventOpt   func(event *ChaincodeEvent) error
)

func WithDefaultSigner(defaultSigner msp.SigningIdentity) Opt {
	return func(opts *Opts) {
		opts.Context = append(opts.Context, func(ctx context.Context) context.Context {
			return ContextWithDefaultSigner(ctx, defaultSigner)
		})
	}
}

func WithTransientValue(key string, value []byte) Opt {
	return func(o *Opts) {
		o.Context = append(o.Context, func(ctx context.Context) context.Context {
			return ContextWithTransientValue(ctx, key, value)
		})
	}
}

func WithEncryption(encKey []byte) Opt {
	return func(o *Opts) {
		WithTransientValue(encryption.TransientMapKey, encKey)(o)
		WithArgsEncryption(encKey)(o)
		WithInvokePayloadDecryption(encKey)(o)
		WithEventDecryption(encKey)(o)
	}
}

func WithArgsEncryption(encKey []byte) Opt {
	return func(o *Opts) {
		o.Input = append(o.Input, func(action Action, ccInput *ChaincodeInput) (err error) {
			ccInput.Args, err = encryption.EncryptArgsBytes(encKey, ccInput.Args)
			return err
		})
	}
}

func WithInvokePayloadDecryption(encKey []byte) Opt {
	return func(o *Opts) {
		o.Output = append(o.Output, func(action Action, r *peer.Response) (err error) {
			if action != Invoke {
				return nil
			}
			r.Payload, err = encryption.Decrypt(encKey, r.Payload)
			return err
		})
	}
}

func WithEventDecryption(encKey []byte) Opt {
	return func(o *Opts) {
		o.Event = append(o.Event, func(e *ChaincodeEvent) error {
			de, err := encryption.DecryptEvent(encKey, e.Event)
			if err != nil {
				return err
			}

			e.Event = de
			return nil
		})
	}
}

func WithEventResolver(resolver mapping.EventResolver) Opt {
	return func(o *Opts) {
		o.Event = append(o.Event, func(e *ChaincodeEvent) error {

			eventPayload, err := resolver.Resolve(e.Event.EventName, e.Event.Payload)
			if err != nil {
				return err
			}

			bb, err := convert.ToBytes(eventPayload)
			if err != nil {
				return err
			}

			e.Event.Payload = bb
			return nil
		})
	}
}
