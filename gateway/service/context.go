package service

import (
	"context"

	"github.com/hyperledger/fabric/msp"
)

const CtxSignerKey = `SigningIdentity`

func ContextWithDefaultSigner(ctx context.Context, defaultSigner msp.SigningIdentity) context.Context {
	if _, err := SignerFromContext(ctx); err != nil {
		return ContextWithSigner(ctx, defaultSigner)
	} else {
		return ctx
	}
}

func ContextWithSigner(ctx context.Context, signer msp.SigningIdentity) context.Context {
	return context.WithValue(ctx, CtxSignerKey, signer)
}

func SignerFromContext(ctx context.Context) (msp.SigningIdentity, error) {
	if signer, ok := ctx.Value(CtxSignerKey).(msp.SigningIdentity); !ok {
		return nil, ErrSignerNotDefinedInContext
	} else {
		return signer, nil
	}
}
