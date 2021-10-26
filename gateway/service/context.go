package service

import (
	"context"

	"github.com/hyperledger/fabric/msp"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	CtxSignerKey   = contextKey(`SigningIdentity`)
	CtxTxWaiterKey = contextKey(`TxWaiter`)
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

func ContextWithTxWaiter(ctx context.Context, txWaiterType string) context.Context {
	return context.WithValue(ctx, CtxTxWaiterKey, txWaiterType)
}

// TxWaiterFromContext - fetch 'txWaiterType' param which identify transaction waiting policy
// what params you'll have depends on your implementation
// for example, in hlf-sdk:
// available: 'self'(wait for one peer of endorser org), 'all'(wait for each organizations from endorsement policy)
// default is 'self'(even if you pass empty string)
func TxWaiterFromContext(ctx context.Context) string {
	txWaiter, _ := ctx.Value(CtxTxWaiterKey).(string)
	return txWaiter
}
