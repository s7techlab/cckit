package gateway

import (
	"context"
)

const CtxTransientKey = `TransientMap`

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
