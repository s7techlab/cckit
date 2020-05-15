package service

import (
	"context"

	"github.com/hyperledger/fabric/msp"
	"github.com/s7techlab/hlf-sdk-go/api"
)

const (
	CtxSignerKey   = `SigningIdentity`
	CtxDoOptionKey = `SdkDoOption`
)

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

func ContextWithDefaultDoOption(ctx context.Context, defaultDoOpts ...api.DoOption) context.Context {
	if opts := DoOptionFromContext(ctx); len(opts) == 0 {
		return ContextWithDoOption(ctx, defaultDoOpts...)
	} else {
		return ctx
	}
}

func ContextWithDoOption(ctx context.Context, doOpts ...api.DoOption) context.Context {
	return context.WithValue(ctx, CtxDoOptionKey, doOpts)
}

func DoOptionFromContext(ctx context.Context) []api.DoOption {
	doOpts, ok := ctx.Value(CtxDoOptionKey).([]api.DoOption)
	if !ok {
		doOpts = []api.DoOption{}
	}
	return doOpts
}
