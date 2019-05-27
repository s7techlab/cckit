package encryption

import (
	"github.com/s7techlab/cckit/router"
)

func EncryptInvokeResponse() router.MiddlewareFunc {
	return func(pre router.HandlerFunc, pos ...int) router.HandlerFunc {
		return func(ctx router.Context) (interface{}, error) {
			res, err := pre(ctx)

			if err != nil || ctx.Handler().Type != router.MethodInvoke {
				return res, err
			}

			return EncryptWithTransientKey(ctx, res)
		}
	}
}
