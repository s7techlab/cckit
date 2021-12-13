package gateway

import (
	"context"

	"github.com/hyperledger/fabric/msp"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	CtxTransientKey = contextKey(`TransientMap`)
	CtxSignerKey    = contextKey(`SigningIdentity`)
	CtxTxWaiterKey  = contextKey(`TxWaiter`)
)

func ContextWithTransientMap(ctx context.Context, transient map[string][]byte) context.Context {
	return context.WithValue(ctx, CtxTransientKey, transient)
}

func ContextWithTransientValue(ctx context.Context, key string, value []byte) context.Context {
	transient, ok := ctx.Value(CtxTransientKey).(map[string][]byte)
	if !ok {
		transient = make(map[string][]byte)
	}
	transient[key] = value
	return context.WithValue(ctx, CtxTransientKey, transient)
}

func TransientFromContext(ctx context.Context) (map[string][]byte, error) {
	if transient, ok := ctx.Value(CtxTransientKey).(map[string][]byte); !ok {
		return nil, nil
	} else {
		return transient, nil
	}
}

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
