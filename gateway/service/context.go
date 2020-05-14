package service

import (
	"context"
	"github.com/s7techlab/hlf-sdk-go/api"

	"github.com/hyperledger/fabric/msp"
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
	if opts, _ := DoOptionFromContext(ctx); len(opts) == 0 {
		return ContextWithDoOption(ctx, defaultDoOpts...)
	} else {
		return ctx
	}
}

func ContextWithDoOption(ctx context.Context, doOpts ...api.DoOption) context.Context {
	return context.WithValue(ctx, CtxDoOptionKey, doOpts)
}

func DoOptionFromContext(ctx context.Context) ([]api.DoOption, error) {
	doOpts := []api.DoOption{}
	doOpts = ctx.Value(CtxDoOptionKey).([]api.DoOption)
	return doOpts, nil
}
