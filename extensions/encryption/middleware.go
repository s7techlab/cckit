package encryption

import (
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/response"
	"github.com/s7techlab/cckit/router"
)

func ArgsDecryptIfKeyProvided(next router.ContextHandlerFunc, pos ...int) router.ContextHandlerFunc {
	return func(c router.Context) peer.Response {
		transient, err := c.Stub().GetTransient()
		if err != nil {
			return response.Error(err)
		}

		key, ok := transient[TransientMapKey]
		// no key provided
		if !ok {
			c.Logger().Debug(`no decrypt key provided`)
			return next(c)
		}

		args, err := DecryptArgs(key, c.Stub().GetArgs())
		if err != nil {
			return response.Error(errors.Wrap(err, `decrypt args error`))
		}

		return next(c.ReplaceArgs(args))
	}
}
