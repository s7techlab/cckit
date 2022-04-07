package encryption

import (
	"github.com/s7techlab/cckit/router"
)

func encryptIfEncodeKeyExists(encryptionRequired bool) router.MiddlewareFunc {
	return func(pre router.HandlerFunc, pos ...int) router.HandlerFunc {
		return func(ctx router.Context) (interface{}, error) {
			res, err := pre(ctx)

			if err != nil || ctx.Handler().Type != router.MethodInvoke {
				return res, err
			}
			if _, err = KeyFromTransient(ctx); err != nil {
				if err == ErrKeyNotDefinedInTransientMap && !encryptionRequired {
					return res, nil
				}
				return nil, err
			}

			return EncryptWithTransientKey(ctx, res)
		}
	}
}

// EncryptInvokeResponse returns middleware function for encrypt chaincode invocation response
func EncryptInvokeResponse() router.MiddlewareFunc {
	return encryptIfEncodeKeyExists(true)
}

// EncryptInvokeResponseIfKeyProvided returns middleware function for encrypt chaincode invocation response if
// encryption key provided in transient map
func EncryptInvokeResponseIfKeyProvided() router.MiddlewareFunc {
	return encryptIfEncodeKeyExists(false)
}
