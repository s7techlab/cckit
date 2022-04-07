package encryption

import (
	"github.com/s7techlab/cckit/router"
)

func responseEncryptor(encryptionRequired bool) router.MiddlewareFunc {
	return func(pre router.HandlerFunc, pos ...int) router.HandlerFunc {
		return func(ctx router.Context) (interface{}, error) {
			res, err := pre(ctx)

			if err != nil || ctx.Handler().Type != router.MethodInvoke {
				return res, err
			}
			var (
				enc []byte
			)
			if enc, err = EncryptWithTransientKey(ctx, res); err == ErrKeyNotDefinedInTransientMap && !encryptionRequired {
				return res, nil
			}

			return enc, err
		}
	}
}

func EncryptInvokeResponse() router.MiddlewareFunc {
	return responseEncryptor(true)
}

func EncryptInvokeResponseIfKeyProvided() router.MiddlewareFunc {
	return responseEncryptor(false)
}
