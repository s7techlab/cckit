package gateway

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/msp"

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

	InstanceOpts struct {
		Locator *ChaincodeLocator
		Opts    Opts
	}

	OptFunc func(*Opts)
	// Deprecated: use OptFunc
	Opt = OptFunc

	ContextOpt func(ctx context.Context) context.Context
	InputOpt   func(input *ChaincodeInput) error
	OutputOpt  func(action InvocationType, response *peer.Response) error
	EventOpt   func(event *ChaincodeEvent) error
)

func WithDefaultSigner(defaultSigner msp.SigningIdentity) OptFunc {
	return func(opts *Opts) {
		opts.Context = append(opts.Context, func(ctx context.Context) context.Context {
			return ContextWithDefaultSigner(ctx, defaultSigner)
		})
	}
}

func WithDefaultTransientMapValue(key string, value []byte) OptFunc {
	return func(o *Opts) {
		o.Input = append(o.Input, func(input *ChaincodeInput) error {
			if input.Transient == nil {
				input.Transient = make(map[string][]byte)
			}
			if _, exists := input.Transient[key]; !exists {
				input.Transient[key] = value
			}
			return nil
		})
	}
}

func WithEncryption(encKey []byte) OptFunc {
	return func(o *Opts) {
		WithDefaultTransientMapValue(encryption.TransientMapKey, encKey)(o)
		WithArgsEncryption(encKey)(o)
		WithInvokePayloadDecryption(encKey)(o)
		WithEventDecryption(encKey)(o)
	}
}

func WithArgsEncryption(encKey []byte) OptFunc {
	return func(o *Opts) {
		o.Input = append(o.Input, func(ccInput *ChaincodeInput) (err error) {
			ccInput.Args, err = encryption.EncryptArgsBytes(encKey, ccInput.Args)
			return err
		})
	}
}

func WithInvokePayloadDecryption(encKey []byte) OptFunc {
	return func(o *Opts) {
		o.Output = append(o.Output, func(action InvocationType, r *peer.Response) (err error) {
			if action != InvocationType_INVOCATION_TYPE_INVOKE {
				return nil
			}
			r.Payload, err = encryption.Decrypt(encKey, r.Payload)
			if err != nil {
				return fmt.Errorf(`decrypt invoke payload: %w`, err)
			}
			return nil
		})
	}
}

func WithEventDecryption(encKey []byte) OptFunc {
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

func WithEventResolver(resolver mapping.EventResolver) OptFunc {
	return func(o *Opts) {
		o.Event = append(o.Event, func(e *ChaincodeEvent) error {
			eventPayload, err := resolver.Resolve(e.Event.EventName, e.Event.Payload)
			if err != nil {
				return err
			}

			bb, err := (&jsonpb.Marshaler{EmitDefaults: true, OrigName: true}).MarshalToString(eventPayload.(proto.Message))
			if err != nil {
				return err
			}

			e.Payload = &RawJson{Value: []byte(bb)}
			return nil
		})
	}
}
